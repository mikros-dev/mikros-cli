package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/cmd/setup"
)

var (
	setupConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Create and install default settings",
		Long: `config creates and installs all default settings into the
CLI TOML file, located in $HOME/.mikros/config.toml`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := setup.CreateDefaultSettings(); err != nil {
				fmt.Println("setup:", err)
				return
			}
		},
	}
)

func setupConfigCmdInit() {
	setupCmd.AddCommand(setupConfigCmd)
}
