package cli

import (
	"os"

	"github.com/fatih/color"
	"github.com/gptscript-ai/cmd"
	"github.com/gptscript-ai/gptscript/pkg/env"
	"github.com/orion101-ai/orion101/apiclient"
	"github.com/orion101-ai/orion101/logger"
	"github.com/orion101-ai/orion101/pkg/cli/internal"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type Orion101 struct {
	Debug  bool `usage:"Enable debug logging"`
	Client *apiclient.Client
}

func (a *Orion101) PersistentPre(cmd *cobra.Command, args []string) error {
	if os.Getenv("NO_COLOR") != "" || !term.IsTerminal(int(os.Stdout.Fd())) {
		color.NoColor = true
	}

	if a.Debug {
		logger.SetDebug()
	}

	if a.Client.Token == "" {
		a.Client = a.Client.WithTokenFetcher(internal.Token)
	}

	return nil
}

func New() *cobra.Command {
	root := &Orion101{
		Client: &apiclient.Client{
			BaseURL: env.VarOrDefault("ORION101_BASE_URL", "http://localhost:8080/api"),
			Token:   os.Getenv("ORION101_TOKEN"),
		},
	}
	return cmd.Command(root,
		&Create{root: root},
		&Agents{root: root},
		cmd.Command(&Workflows{root: root},
			&WorkflowAuth{root: root}),
		&Edit{root: root},
		&Update{root: root},
		&Delete{root: root},
		&Invoke{root: root},
		cmd.Command(&Threads{root: root}, &ThreadPrint{root: root}),
		cmd.Command(&Credentials{root: root}, &CredentialsDelete{root: root}),
		cmd.Command(&Runs{root: root}, &Debug{root: root}, &RunPrint{root: root}),
		cmd.Command(&Tools{root: root},
			&ToolUnregister{root: root},
			&ToolRegister{root: root},
			&ToolUpdate{root: root}),
		&Webhooks{root: root},
		&Server{},
		&Version{},
	)
}

func (a *Orion101) Run(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
