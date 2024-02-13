package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	serviceCmd = &cobra.Command{
		Use:     "service",
		Aliases: []string{"svc"},
		Short:   "Handles service tasks.",
		Long:    `service command helps creating or editing services locally.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				fmt.Println("service:", err)
				return
			}
		},
	}
)

func serviceCmdInit() {
	serviceInitCmdInit()
	rootCmd.AddCommand(serviceCmd)
}
