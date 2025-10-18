package commands

import (
	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

// EntryPoint initializes and returns the root command for the application,
// configured with its subcommands.
func EntryPoint(cfg *settings.Settings) *cobra.Command {
	root := rootCmd()

	// Configure commands
	root.AddCommand(configCmd())
	root.AddCommand(newCmd(cfg))
	root.AddCommand(lintCmd())

	return root
}
