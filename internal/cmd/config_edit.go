package cmd

import (
	"fmt"

	"github.com/mikros-dev/mikros-cli/internal/config/edit"
	"github.com/spf13/cobra"
)

var (
	configEditCmd = &cobra.Command{
		Use:   "edit",
		Short: "Edit the configuration file",
		Long: `edit command opens (or creates if it does not exist) the
mikros CLI settings file into a form allowing it to be
customized.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := edit.New(); err != nil {
				fmt.Println("config: ", err)
				return
			}
		},
	}
)

func configEditCmdInit() {
	configCmd.AddCommand(configEditCmd)
}
