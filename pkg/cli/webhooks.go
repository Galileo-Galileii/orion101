package cli

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/spf13/cobra"
)

type Webhooks struct {
	root   *Orion101
	Quiet  bool   `usage:"Only print IDs of agents" short:"q"`
	Wide   bool   `usage:"Print more information" short:"w"`
	Output string `usage:"Output format (table, json, yaml)" short:"o" default:"table"`
}

func (l *Webhooks) Customize(cmd *cobra.Command) {
	cmd.Aliases = []string{"webhook", "wh"}
}

func (l *Webhooks) Run(cmd *cobra.Command, args []string) error {
	whs, err := l.root.Client.ListWebhooks(cmd.Context())
	if err != nil {
		return err
	}

	if ok, err := output(l.Output, whs); ok || err != nil {
		return err
	}

	if l.Quiet {
		for _, webhook := range whs.Items {
			fmt.Println(webhook.ID)
		}
		return nil
	}

	w := newTable("ID", "NAME", "DESCRIPTION", "WORKFLOW", "LASTRUN", "CREATED")
	for _, wh := range whs.Items {
		w.WriteRow(wh.ID, wh.Name, truncate(wh.Description, l.Wide), wh.Workflow,
			humanize.Time(wh.LastSuccessfulRunCompleted.GetTime()),
			humanize.Time(wh.Created.Time))
	}

	return w.Err()
}