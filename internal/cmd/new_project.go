package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mikros-dev/mikros-cli/internal/cmd/new/project"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

var (
	newProjectCmd = &cobra.Command{
		Use:   "project",
		Short: "Create a new mikros project",
		Long:  "project is a command to help creating mikros project",
	}
)

func newProjectCmdInit(cfg *settings.Settings) {
	newProjectCmd.Run = func(cmd *cobra.Command, args []string) {
		if err := project.New(cfg); err != nil {
			fmt.Println("project:", err)
			return
		}

		fmt.Printf("\n✅ Project successfully created\n\n")
		fmt.Println("In order to start, execute the following command inside the new project directory:")
		fmt.Printf("\n$ make setup\n\n")
	}

	newCmd.AddCommand(newProjectCmd)
}
