package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new mikros project",
		Long:  "new helps creating a new mikros project",
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Println("new:", err)
				return
			}
		},
	}
)

func newCmdInit(cfg *settings.Settings) {
	newServiceCmdInit(cfg)
	newProjectCmdInit(cfg)
	rootCmd.AddCommand(newCmd)
}
