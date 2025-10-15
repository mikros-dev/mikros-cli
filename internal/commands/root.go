package commands

import (
	"github.com/spf13/cobra"
)

func rootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mikros",
		Short: "A \"swiss army knife\" for dealing with mikros framework tasks.",
		Long: `mikros is a command to help the developer use the mikros
framework to create new services.`,
	}
}
