package server

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/orion101-ai/orion101/pkg/api"
	"github.com/orion101-ai/orion101/pkg/gateway/context"
	"github.com/orion101-ai/orion101/pkg/gateway/log"
)

func apply(h api.HandlerFunc, m ...api.Middleware) api.HandlerFunc {
	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}
	return h
}

func contentType(contentTypes ...string) api.Middleware {
	return func(h api.HandlerFunc) api.HandlerFunc {
		return func(apiContext api.Context) error {
			for _, ct := range contentTypes {
				apiContext.ResponseWriter.Header().Add("Content-Type", ct)
			}
			return h(apiContext)
		}
	}
}

func logRequest(h api.HandlerFunc) api.HandlerFunc {
	return func(apiContext api.Context) (err error) {
		l := context.GetLogger(apiContext.Context())
		defer func() {
			l.DebugContext(apiContext.Context(), "Handled request", "method", apiContext.Method, "path", apiContext.URL.Path)
			if recErr := recover(); recErr != nil {
				l.ErrorContext(apiContext.Context(), "Panic", "error", err, "stack", string(debug.Stack()))
				err = fmt.Errorf("encountered an unexpected error")
			}
		}()

		l.DebugContext(apiContext.Context(), "Handling request", "method", apiContext.Method, "path", apiContext.URL.Path)
		return h(apiContext)
	}
}

func addRequestID(next api.HandlerFunc) api.HandlerFunc {
	return func(apiContext api.Context) error {
		apiContext.Request = apiContext.Request.WithContext(context.WithNewRequestID(apiContext.Request.Context()))
		return next(apiContext)
	}
}

func addLogger(next api.HandlerFunc) api.HandlerFunc {
	return func(apiContext api.Context) error {
		logger := log.NewWithID(context.GetRequestID(apiContext.Request.Context()))
		if apiContext.User != nil {
			logger = logger.With("username", apiContext.User.GetName())
		}
		apiContext.Request = apiContext.Request.WithContext(context.WithLogger(
			apiContext.Request.Context(),
			logger,
		))
		return next(apiContext)
	}
}

func httpToApiHandlerFunc(handler http.Handler) api.HandlerFunc {
	return func(apiContext api.Context) error {
		handler.ServeHTTP(apiContext.ResponseWriter, apiContext.Request)
		return nil
	}
}