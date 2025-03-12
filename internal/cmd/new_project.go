package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

var (
	newProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Create a new mikros project",
		Long:  "project is a command to help creating mikros project",
		Run:   func(cmd *cobra.Command, args []string) {},
	}
)

func newProjectCmdInit(_ *settings.Settings) {
	newCmd.AddCommand(newProjectCmd)
}
