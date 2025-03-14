package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "mikros",
		Short: "A \"swiss army knife\" for dealing with mikros framework tasks.",
		Long: `mikros is a command to help the developer use the mikros
framework to create new services.`,
	}
)

// Execute puts the CLI to execute.
func Execute() {
	cfg, err := settings.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	loadCommands(cfg)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loadCommands is where all CLI options are loaded and prepared to be
// executed.
func loadCommands(cfg *settings.Settings) {
	newCmdInit(cfg)
	configCmdInit()
}
