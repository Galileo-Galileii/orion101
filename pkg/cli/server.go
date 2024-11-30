package cli

import (
	"github.com/orion101-ai/orion101/pkg/server"
	"github.com/orion101-ai/orion101/pkg/services"
	"github.com/spf13/cobra"
)

type Server struct {
	services.Config
}

func (s *Server) Run(cmd *cobra.Command, args []string) error {
	return server.Run(cmd.Context(), s.Config)
}
