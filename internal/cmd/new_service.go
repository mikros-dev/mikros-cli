package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikros-dev/mikros-cli/internal/cmd/new/service"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

var (
	newServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "Create a new mikros service",
		Long:  "service is a helper command to create a new mikros service",
	}
)

func newServiceCmdInit(cfg *settings.Settings) {
	// path option
	newServiceCmd.Flags().String("path", "", "Sets the output path name (default: cwd).")
	_ = viper.BindPFlag("service-path", newServiceCmd.Flags().Lookup("path"))

	// proto file option
	newServiceCmd.Flags().String("proto", "", "Uses an _api.proto file as source for the service API.")
	_ = viper.BindPFlag("service-proto", newServiceCmd.Flags().Lookup("proto"))

	// sets the function handler
	newServiceCmd.Run = func(cmd *cobra.Command, args []string) {
		initOptions := &service.InitOptions{
			Path:          viper.GetString("service-path"),
			ProtoFilename: viper.GetString("service-proto"),
		}

		if err := service.Init(cfg, initOptions); err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("\n✅ Service successfully created\n")
	}

	newCmd.AddCommand(newServiceCmd)
}
