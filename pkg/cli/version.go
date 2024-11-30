package cli

import (
	"fmt"

	"github.com/orion101-ai/orion101/pkg/version"
	"github.com/spf13/cobra"
)

type Version struct {
	root *Orion101
}

func (l *Version) Run(cmd *cobra.Command, args []string) error {
	fmt.Println("Version: ", version.Get())
	return nil
}
