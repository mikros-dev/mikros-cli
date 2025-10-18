package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikros-dev/mikros-cli/internal/lint"
)

func lintCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Analyze the code with the project's rules",
		Long: `Runs the revive linter with the project's default configuration.
Use it to validate style and potential issues before a commit.

It is idempotent and does not modify files by default. Use flags
to customize paths and format.

Examples:
 # Run lint in the current directory
 $ mikros lint

 # Run lint excluding a specific file
 $ mikros lint --exclude foo.go

 # Run lint excluding all files in a directory
 $ mikros lint --exclude foo/...

 # Run lint excluding more than one file/directory
 $ mikros lint --exclude foo/...,bar/...,file.go
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return lint.Run(cmd.Context(), lint.Options{
				Debug:   viper.GetBool("lint.debug"),
				Format:  viper.GetString("lint.format"),
				Config:  viper.GetString("lint.config"),
				Path:    viper.GetString("lint.path"),
				Exclude: strings.Split(viper.GetString("lint.exclude"), ","),
			})
		},
	}

	addLintCommandFlags(cmd)

	return cmd
}

func addLintCommandFlags(cmd *cobra.Command) {
	cmd.Flags().String("format", "friendly", "Output format (e.g., friendly, json)")
	cmd.Flags().String("path", "./...", "Path to analyze")
	cmd.Flags().String("exclude", "", "Directories to exclude from validation")
	cmd.Flags().Bool("debug", false, "Enable debug mode")

	_ = viper.BindPFlag("lint.format", cmd.Flags().Lookup("format"))
	_ = viper.BindPFlag("lint.path", cmd.Flags().Lookup("path"))
	_ = viper.BindPFlag("lint.exclude", cmd.Flags().Lookup("exclude"))
	_ = viper.BindPFlag("lint.debug", cmd.Flags().Lookup("debug"))
}
