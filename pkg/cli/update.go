package cli

import (
	"fmt"
	"os"

	"github.com/orion101-ai/orion101/apiclient/types"
	"github.com/orion101-ai/orion101/pkg/system"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type Update struct {
	root  *Orion101
	Quiet bool `usage:"Only print IDs of updated agent/workflow" short:"q"`
}

func (l *Update) Customize(cmd *cobra.Command) {
	cmd.Use = "update [flags] [ID] [MANIFEST_FILE]"
	cmd.Args = cobra.ExactArgs(2)
}

func (l *Update) Run(cmd *cobra.Command, args []string) error {
	id := args[0]
	data, err := os.ReadFile(args[1])
	if err != nil {
		return err
	}

	var newManifest types.WorkflowManifest
	if err := yaml.Unmarshal(data, &newManifest); err != nil {
		return err
	}

	if system.IsWorkflowID(id) {
		wf, err := l.root.Client.UpdateWorkflow(cmd.Context(), id, newManifest)
		if err != nil {
			return err
		}
		if l.Quiet {
			fmt.Println(wf.ID)
			return nil
		}
		fmt.Printf("Workflow updated: %s\n", wf.ID)
		return nil
	}

	agent, err := l.root.Client.UpdateAgent(cmd.Context(), id, newManifest.AgentManifest)
	if err != nil {
		return err
	}
	if l.Quiet {
		fmt.Println(agent.ID)
		return nil
	}
	fmt.Printf("Agent updated: %s\n", agent.ID)
	return nil
}