package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/cmd/config"
)

var (
	configSetupCmd = &cobra.Command{
		Use:   "setup",
		Short: "Create and install default settings",
		Long: `setup creates and installs all default settings into the
CLI TOML file, located in $HOME/.mikros/config.toml`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.CreateDefaultSettings(); err != nil {
				fmt.Println("config:", err)
				return
			}
		},
	}
)

func configSetupCmdInit() {
	configCmd.AddCommand(configSetupCmd)
}
