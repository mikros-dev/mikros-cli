package cmd

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikros-dev/mikros-cli/internal/scaffold/protobuf"
	protobuf_repository "github.com/mikros-dev/mikros-cli/internal/scaffold/repository/protobuf"
	service_repository "github.com/mikros-dev/mikros-cli/internal/scaffold/repository/service"
	"github.com/mikros-dev/mikros-cli/internal/scaffold/service"
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func cmd(cfg *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new mikros project",
		Long:  "new helps creating a new mikros project",
		RunE: func(cmd *cobra.Command, args []string) error {
			selected, err := runNewProjectForm(cfg)
			if err != nil {
				return err
			}

			switch selected {
			case "protobuf-monorepo":
				return newProtobufRepository(cfg)

			case "services-monorepo":
				return newServiceRepository(cfg)

			case "protobuf-module":
				return newProtobufModule(cfg)

			case "service-template":
				return newServiceTemplate(cfg)

			case "quit":
				// Just quits
				return nil
			}

			return nil
		},
	}

	setNewCmdFlags(cmd)

	return cmd
}

func newProtobufRepository(cfg *settings.Settings) error {
	options := &protobuf_repository.NewOptions{
		NoVCS:   viper.GetBool("project-no-vcs"),
		Path:    viper.GetString("project-path"),
		Profile: viper.GetString("project-profile"),
	}

	if err := protobuf_repository.New(cfg, options); err != nil {
		return err
	}

	fmt.Printf("\n✅ Project successfully created\n\n")
	fmt.Println("In order to start, execute the following command inside the new project directory:")
	fmt.Printf("\n$ make setup\n\n")

	return nil
}

func newServiceRepository(cfg *settings.Settings) error {
	options := &service_repository.NewOptions{
		NoVCS: viper.GetBool("project-no-vcs"),
		Path:  viper.GetString("project-path"),
	}

	if err := service_repository.New(cfg, options); err != nil {
		return err
	}

	fmt.Printf("\n✅ Project successfully created\n\n")
	return nil
}

func newProtobufModule(cfg *settings.Settings) error {
	options := &protobuf.NewOptions{
		Profile: viper.GetString("project-profile"),
	}

	return protobuf.New(cfg, options)
}

func newServiceTemplate(cfg *settings.Settings) error {
	options := &service.NewOptions{
		Path:          viper.GetString("project-path"),
		ProtoFilename: viper.GetString("project-proto"),
	}

	if err := service.New(cfg, options); err != nil {
		return err
	}

	fmt.Printf("\n✅ Service successfully created\n")
	return nil
}

func setNewCmdFlags(cmd *cobra.Command) {
	// path option
	cmd.Flags().String("path", "", "Sets the output path name (default cwd).")
	_ = viper.BindPFlag("project-path", cmd.Flags().Lookup("path"))

	// proto file option
	cmd.Flags().String("proto", "", "Uses an _api.proto file as source for the service API.")
	_ = viper.BindPFlag("project-proto", cmd.Flags().Lookup("proto"))

	// no-vcs option
	cmd.Flags().Bool("no-vcs", false, "Disables creating projects with VCS support (default true).")
	_ = viper.BindPFlag("project-no-vcs", cmd.Flags().Lookup("no-vcs"))

	// profile option
	cmd.Flags().String("profile", "default", "Sets the profile to use.")
	_ = viper.BindPFlag("project-profile", cmd.Flags().Lookup("profile"))
}

func runNewProjectForm(cfg *settings.Settings) (string, error) {
	var selectedProject string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a project to create or Quit to exit the application").
				Options(
					huh.NewOption("Protobuf monorepo", "protobuf-monorepo"),
					huh.NewOption("Services monorepo", "services-monorepo"),
					huh.NewOption("Protobuf module file(s)", "protobuf-module"),
					huh.NewOption("Single service template", "service-template"),
					huh.NewOption("Quit", "quit"),
				).
				Value(&selectedProject),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return selectedProject, nil
}
