package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	configCmd = &cobra.Command{
		Use:   "config",
		Short: "Set up mikros related requirements",
		Long: `config helps installing and adjusting mikros related requirements
inside the system.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Println("config:", err)
				return
			}
		},
	}
)

func configCmdInit() {
	configGenerateCmdInit()
	configEditCmdInit()
	rootCmd.AddCommand(configCmd)
}
