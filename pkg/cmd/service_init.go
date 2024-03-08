package cmd

import (
	"embed"
	"fmt"
	"slices"

	"github.com/somatech1/mikros/components/plugin"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/somatech1/mikros-cli/internal/cmd/service"
	"github.com/somatech1/mikros-cli/pkg/templates"
)

type serviceInitCmdOptions struct {
	DisablePersistentFlags bool
	Features               *plugin.FeatureSet
	Services               *plugin.ServiceSet
	AdditionalTemplates    *ServiceTemplateFile
}

type ServiceTemplateFile struct {
	Files     embed.FS
	Templates []templates.TemplateFile
	Api       map[string]interface{}
}

var (
	serviceInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes a new service.",
		Long: `init helps creating a new service template by creating
some of its source files and its go.mod.`,
	}
)

func serviceInitCmdInit(options *serviceInitCmdOptions) {
	// path option
	serviceInitCmd.Flags().String("path", "", "Sets the output path name (default: cwd).")
	_ = viper.BindPFlag("init-path", serviceInitCmd.Flags().Lookup("path"))

	// proto file option
	serviceInitCmd.Flags().String("proto", "", "Uses an _api.proto file as source for the service API.")
	_ = viper.BindPFlag("init-proto", serviceInitCmd.Flags().Lookup("proto"))

	serviceInitCmd.Run = func(cmd *cobra.Command, args []string) {
		initOptions := &service.InitOptions{
			Path:          viper.GetString("init-path"),
			ProtoFilename: viper.GetString("init-proto"),
		}

		if options != nil {
			initOptions.Features = options.Features
			initOptions.Services = options.Services

			if options.AdditionalTemplates != nil {
				initOptions.ExternalTemplates = &service.TemplateFileOptions{
					Files:     options.AdditionalTemplates.Files,
					Templates: options.AdditionalTemplates.Templates,
					Api:       options.AdditionalTemplates.Api,
				}
			}
		}

		if err := service.Init(initOptions); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("\n✅ Service successfully created\n")
	}

	if options.DisablePersistentFlags {
		serviceInitCmd.SetHelpFunc(func(command *cobra.Command, i []string) {
			disableServiceGlobalFlags()
			command.Parent().HelpFunc()(command, i)
		})
	}

	serviceCmd.AddCommand(serviceInitCmd)
}

func disableServiceGlobalFlags() {
	flagsToHide := []string{}
	serviceCmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		flagsToHide = append(flagsToHide, flag.Name)
	})

	markGlobalFlagsHidden(serviceCmd, flagsToHide...)
}

func markGlobalFlagsHidden(command *cobra.Command, flags ...string) {
	command.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		if slices.Contains(flags, flag.Name) {
			flag.Hidden = true
		}
	})
}
