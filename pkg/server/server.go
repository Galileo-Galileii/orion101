package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/orion101-ai/orion101/logger"
	"github.com/orion101-ai/orion101/pkg/api/router"
	"github.com/orion101-ai/orion101/pkg/controller"
	"github.com/orion101-ai/orion101/pkg/services"
	"github.com/rs/cors"
)

var log = logger.Package()

func Run(ctx context.Context, c services.Config) error {
	svcs, err := services.New(ctx, c)
	if err != nil {
		return err
	}

	go func() {
		c, err := controller.New(svcs)
		if err != nil {
			log.Fatalf("Failed to start controller: %v", err)
		}
		if err := c.Start(ctx); err != nil {
			log.Fatalf("Failed to start controller: %v", err)
		}
		if err := c.PostStart(ctx); err != nil {
			log.Fatalf("Failed to post start controller: %v", err)
		}
	}()

	handler, err := router.Router(svcs)
	if err != nil {
		return err
	}

	context.AfterFunc(ctx, func() {
		log.Fatalf("Interrupted, exiting")
	})

	if c.DevMode && c.AllowedOrigin == "" {
		c.AllowedOrigin = "*"
	}

	address := fmt.Sprintf("0.0.0.0:%d", c.HTTPListenPort)
	log.Infof("Starting server on %s", address)
	allowEverything := cors.New(cors.Options{
		AllowedOrigins: []string{c.AllowedOrigin},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"*"},
	})
	return http.ListenAndServe(address, allowEverything.Handler(handler))
}