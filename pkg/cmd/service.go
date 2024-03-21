package cmd

import (
	"embed"
	"fmt"

	"github.com/somatech1/mikros/components/plugin"
	"github.com/spf13/cobra"

	"github.com/somatech1/mikros-cli/pkg/templates"
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

// ServiceCmdOptions gathers options to set the behavior of the 'service'
// command when integrating it into another CLI.
type ServiceCmdOptions struct {
	// DisablePersistentFlagsOnInit if true removes custom persistent flags
	// from the 'init' subcommand, since they can't be used by it.
	DisablePersistentFlagsOnInit bool

	// Features adds a set of custom features that can be enabled/disabled
	// when creating a new template service.
	Features *plugin.FeatureSet

	// Services adds a set of custom service types that can be selected when
	// creating a new template service.
	Services *plugin.ServiceSet

	// SubCommands adds new custom subcommands for the 'service' command.
	SubCommands []*cobra.Command

	// PersistentFlags adds new custom flags that will persist for all available
	// subcommands.
	PersistentFlags []*Flag

	// AdditionalTemplates adds additional templates that will be generated
	// according selected options when creating a new template service.
	AdditionalTemplates *ServiceTemplateFile
}

// Flag represents a new flag option to be added to the command.
type Flag struct {
	// Name the command line flag name.
	Name string

	// Usage is a short description for the flag.
	Usage string

	// Value holds the default value of the flag. At the moment, only string and
	// bool are supported.
	Value interface{}
}

// ServiceTemplateFile gathers custom template options.
type ServiceTemplateFile struct {
	// Files include all available template files that will compose the generated
	// templates of a new service.
	Files embed.FS

	// Templates gathers detailed information related to all templates being added.
	Templates []templates.TemplateFile

	// Api gathers custom API functions that will be available for all templates.
	Api map[string]interface{}

	// NewServiceArgs allows adding custom content for external service kind when
	// creating the main file of a new template service. It will be available inside
	// the {{.NewServiceArgs}} inside the template.
	//
	// It can be a template string that receives the internal default API and
	// the custom Api as functions available to be used as well as a short data
	// context with the following fields:
	//
	// {{.ServiceName}}: holding the current service name.
	// {{.ServiceType}}: holding the current service type.
	// {{.ServiceTypeCustomAnswers}}: holding custom answers related to the current
	//								  service type.
	NewServiceArgs map[string]string

	// WithExternalFeaturesArg sets the argument of method .WithExternalFeatures()
	// inside the main template.
	WithExternalFeaturesArg string

	// WithExternalServicesArg sets the argument of method .WithExternalServices()
	// inside the main template.
	WithExternalServicesArg string
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

	serviceInitCmdInit(&serviceInitCmdOptions{
		DisablePersistentFlags: options.DisablePersistentFlagsOnInit,
		Features:               options.Features,
		Services:               options.Services,
		AdditionalTemplates:    options.AdditionalTemplates,
	})

	serviceCmd.AddCommand(options.SubCommands...)
	return serviceCmd
}

// CurrentServiceCommand gives access to the current service command without
// any change.
func CurrentServiceCommand() *cobra.Command {
	return serviceInitCmd
}
