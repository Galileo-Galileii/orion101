package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/orion101-ai/orion101/apiclient"
	"github.com/spf13/cobra"
)

type Workflows struct {
	root   *Orion101
	Quiet  bool   `usage:"Only print IDs of agents" short:"q"`
	Wide   bool   `usage:"Print more information" short:"w"`
	Output string `usage:"Output format (table, json, yaml)" short:"o" default:"table"`
}

func (l *Workflows) Customize(cmd *cobra.Command) {
	cmd.Aliases = []string{"workflow", "wf", "w"}
}

func (l *Workflows) Run(cmd *cobra.Command, args []string) error {
	wfs, err := l.root.Client.ListWorkflows(cmd.Context(), apiclient.ListWorkflowsOptions{})
	if err != nil {
		return err
	}

	if ok, err := output(l.Output, wfs); ok || err != nil {
		return err
	}

	if l.Quiet {
		for _, agent := range wfs.Items {
			fmt.Println(agent.ID)
		}
		return nil
	}

	w := newTable("ID", "NAME", "DESCRIPTION", "INVOKE", "CREATED")
	for _, wf := range wfs.Items {
		w.WriteRow(wf.ID, wf.Name, truncate(wf.Description, l.Wide), wf.Links["invoke"], humanize.Time(wf.Created.Time))
	}

	return w.Err()
}