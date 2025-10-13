package commands

import (
	"github.com/spf13/cobra"
)

func configCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Set up mikros related requirements",
		Long: `config helps installing and adjusting mikros related requirements
inside the system.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(configEditCmd())
	cmd.AddCommand(configGenerateCmd())

	return cmd
}
