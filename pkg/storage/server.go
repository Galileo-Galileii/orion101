package storage

import (
	"context"
	"encoding/json"
	"net"
	"net/http"

	"github.com/orion101-ai/kinm/pkg/db"
	mserver "github.com/orion101-ai/kinm/pkg/server"
	"github.com/orion101-ai/orion101/pkg/storage/openapi/generated"
	"github.com/orion101-ai/orion101/pkg/storage/registry/apigroups/agent"
	"github.com/orion101-ai/orion101/pkg/storage/scheme"
	"github.com/orion101-ai/orion101/pkg/storage/services"
	"github.com/orion101-ai/orion101/pkg/version"
	k8sversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client client.WithWatch

func Start(ctx context.Context, config services.Config) (Client, *rest.Config, *db.Factory, error) {
	services, err := services.New(config)
	if err != nil {
		return nil, nil, nil, err
	}

	c, cfg, err := startMinkServer(ctx, config, services)
	if err != nil {
		return nil, nil, nil, err
	}
	return c, cfg, services.DB, nil
}

func startMinkServer(ctx context.Context, config services.Config, services *services.Services) (Client, *rest.Config, error) {
	apiGroups, err := mserver.BuildAPIGroups(services, agent.APIGroup, agent.LeasesAPIGroup)
	if err != nil {
		return nil, nil, err
	}

	var l net.Listener
	if config.StorageListenPort == 0 {
		l, err = net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, nil, err
		}
	}

	minkConfig := &mserver.Config{
		Name:              "Storage Server",
		Version:           version.Get().String(),
		Authenticator:     services.Authn,
		Authorization:     services.Authz,
		HTTPSListenPort:   config.StorageListenPort,
		Listener:          l,
		OpenAPIConfig:     generated.GetOpenAPIDefinitions,
		Scheme:            scheme.Scheme,
		APIGroups:         apiGroups,
		ReadinessCheckers: []healthz.HealthChecker{services.DB},
	}

	//if cfg.AuditLogPolicyFile != "" && cfg.AuditLogPath != "" {
	//	minkConfig.AuditConfig = mserver.NewAuditOptions(cfg.AuditLogPolicyFile, cfg.AuditLogPath)
	//}

	minkConfig.Middleware = []func(http.Handler) http.Handler{
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.URL.Path == "/version" {
					_ = json.NewEncoder(rw).Encode(k8sversion.Info{
						GitVersion: version.Get().String(),
						GitCommit:  version.Get().Commit,
					})
				} else {
					next.ServeHTTP(rw, req)
				}
			})
		},
	}

	minkServer, err := mserver.New(minkConfig)
	if err != nil {
		return nil, nil, err
	}

	_ = minkServer.Handler(ctx)

	cfg := minkServer.Loopback
	c, err := client.NewWithWatch(cfg, client.Options{
		Scheme: scheme.Scheme,
	})
	return c, cfg, err
}