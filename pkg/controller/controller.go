package controller

import (
	"context"
	"fmt"

	"github.com/orion101-ai/nah/pkg/router"
	"github.com/orion101-ai/orion101/pkg/controller/data"
	"github.com/orion101-ai/orion101/pkg/controller/handlers/toolreference"
	"github.com/orion101-ai/orion101/pkg/services"
	// Enable logrus logging in nah
	_ "github.com/orion101-ai/nah/pkg/logrus"
)

type Controller struct {
	router         *router.Router
	services       *services.Services
	toolRefHandler *toolreference.Handler
}

func New(services *services.Services) (*Controller, error) {
	c := &Controller{
		router:   services.Router,
		services: services,
	}

	err := c.setupRoutes()
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Controller) PostStart(ctx context.Context) error {
	if err := data.Data(ctx, c.services.StorageClient); err != nil {
		return fmt.Errorf("failed to apply data: %w", err)
	}
	go c.toolRefHandler.PollRegistry(ctx, c.services.Router.Backend())
	return nil
}

func (c *Controller) Start(ctx context.Context) error {
	if err := c.router.Start(ctx); err != nil {
		return fmt.Errorf("failed to start router: %w", err)
	}
	return nil
}
