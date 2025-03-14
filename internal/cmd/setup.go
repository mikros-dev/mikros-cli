package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	setupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Set up mikros related requirements",
		Long: `setup helps installing and adjusting mikros related requirements
inside the system.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Println("setup:", err)
				return
			}
		},
	}
)

func setupCmdInit() {
	setupConfigCmdInit()
	rootCmd.AddCommand(setupCmd)
}
