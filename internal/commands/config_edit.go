package commands

import (
	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/config/edit"
)

func configEditCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration file",
		Long: `edit command opens (or creates if it does not exist) the
mikros CLI settings file into a form allowing it to be
customized.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return edit.New()
		},
	}
}
