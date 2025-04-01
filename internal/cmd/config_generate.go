package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/cmd/config"
)

var (
	configGenerateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Create and install default settings",
		Long: `generate creates and installs all default settings into the
CLI TOML file, located in $HOME/.mikros/config.toml`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := config.CreateDefaultSettings(); err != nil {
				fmt.Println("config:", err)
				return
			}
		},
	}
)

func configGenerateCmdInit() {
	configCmd.AddCommand(configGenerateCmd)
}
