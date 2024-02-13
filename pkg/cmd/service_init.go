package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/somatech1/mikros-cli/internal/cmd/service"
)

var (
	serviceInitCmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes a new service.",
		Long: `init helps initialize a new service folder by creating
some of its settings file and its go.mod.`,
		Run: func(cmd *cobra.Command, args []string) {
			options := &service.InitOptions{
				Path:          viper.GetString("init-path"),
				ProtoFilename: viper.GetString("init-proto"),
			}

			if err := service.Init(options); err != nil {
				fmt.Println(err.Error())
				return
			}
		},
	}
)

func serviceInitCmdInit() {
	serviceCmd.AddCommand(serviceInitCmd)

	// path option
	serviceInitCmd.Flags().String("path", "", "Sets the output path name (default: cwd).")
	_ = viper.BindPFlag("init-path", serviceInitCmd.Flags().Lookup("path"))

	// proto file option
	serviceInitCmd.Flags().String("proto", "", "Uses an _api.proto file as source for the service API.")
	_ = viper.BindPFlag("init-proto", serviceInitCmd.Flags().Lookup("proto"))
}
