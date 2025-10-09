package cmd

import (
	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/config"
)

func configGenerateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Create and install default settings",
		Long: `generate creates and installs all default settings into the
CLI TOML file, located in $HOME/.mikros/config.toml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return config.CreateDefaultSettings()
		},
	}
}
