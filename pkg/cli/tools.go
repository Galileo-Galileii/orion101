package cli

import (
	"fmt"

	"github.com/orion101-ai/orion101/apiclient"
	"github.com/spf13/cobra"
)

type Tools struct {
	root   *Orion101
	Quiet  bool   `usage:"Only print IDs of tools" short:"q"`
	Output string `usage:"Output format (table, json, yaml)" short:"o" default:"table"`
}

func (l *Tools) Customize(cmd *cobra.Command) {
	cmd.Use = "tools [flags]"
	cmd.Args = cobra.NoArgs
	cmd.Aliases = []string{"tool", "tl"}
}

func (l *Tools) Run(cmd *cobra.Command, args []string) error {
	toolRefs, err := l.root.Client.ListToolReferences(cmd.Context(), apiclient.ListToolReferencesOptions{})
	if err != nil {
		return err
	}

	if ok, err := output(l.Output, toolRefs); ok || err != nil {
		return err
	}

	if l.Quiet {
		for _, toolRef := range toolRefs.Items {
			fmt.Println(toolRef.ID)
		}
		return nil
	}

	w := newTable("ID", "NAME", "REF", "DESCRIPTION", "TYPE")
	for _, toolRef := range toolRefs.Items {
		desc := toolRef.Description
		if toolRef.Error != "" {
			desc = toolRef.Error
		}
		ref := toolRef.Reference
		if toolRef.Builtin {
			ref = "builtin"
		}
		w.WriteRow(toolRef.ID, toolRef.Name, ref, desc, string(toolRef.ToolType))
	}

	return w.Err()
}
