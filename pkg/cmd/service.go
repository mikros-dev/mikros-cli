package cmd

import (
	"fmt"

	"github.com/somatech1/mikros/components/plugin"
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
	serviceInitCmdInit(nil)
	rootCmd.AddCommand(serviceCmd)
}

type ServiceCmdOptions struct {
	DisablePersistentFlagsOnInit bool
	Features                     *plugin.FeatureSet
	Services                     *plugin.ServiceSet
	SubCommands                  []*cobra.Command
	PersistentFlags              []*Flag
	Flags                        []*Flag
}

type Flag struct {
	Name  string
	Usage string
	Value interface{}
}

// ServiceCommand gives access to the 'service' CLI command in order to be used
// by another application allowing the caller to customize its behavior.
func ServiceCommand(options *ServiceCmdOptions) *cobra.Command {
	for _, flag := range options.PersistentFlags {
		switch v := flag.Value.(type) {
		case string:
			serviceCmd.PersistentFlags().String(flag.Name, v, flag.Usage)
		case bool:
			serviceCmd.PersistentFlags().Bool(flag.Name, v, flag.Usage)
		}
	}

	for _, flag := range options.Flags {
		switch v := flag.Value.(type) {
		case string:
			serviceCmd.Flags().String(flag.Name, v, flag.Usage)
		case bool:
			serviceCmd.Flags().Bool(flag.Name, v, flag.Usage)
		}
	}

	serviceInitCmdInit(&serviceInitCmdOptions{
		DisablePersistentFlags: options.DisablePersistentFlagsOnInit,
		Features:               options.Features,
		Services:               options.Services,
	})

	serviceCmd.AddCommand(options.SubCommands...)
	return serviceCmd
}

// CurrentServiceCommand gives access to the current service command without
// any change.
func CurrentServiceCommand() *cobra.Command {
	return serviceInitCmd
}
