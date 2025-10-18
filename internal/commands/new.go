package commands

import (
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mikros-dev/mikros-cli/internal/scaffold/protobuf"
	protobuf_repository "github.com/mikros-dev/mikros-cli/internal/scaffold/repository/protobuf"
	service_repository "github.com/mikros-dev/mikros-cli/internal/scaffold/repository/service"
	"github.com/mikros-dev/mikros-cli/internal/scaffold/service"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func newCmd(cfg *settings.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new mikros project",
		Long:  "new helps creating different mikros projects",
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
		NoVCS:   viper.GetBool("new.no-vcs"),
		Path:    viper.GetString("new.path"),
		Profile: viper.GetString("new.profile"),
	}

	if err := protobuf_repository.New(cfg, options); err != nil {
		return err
	}

	ui.Message(cfg, "New protobuf repository",
		"✅ Project successfully created \n\n"+
			"In order to start, execute the following command inside the new project directory:"+
			"\n\n$ make setup")

	return nil
}

func newServiceRepository(cfg *settings.Settings) error {
	options := &service_repository.NewOptions{
		NoVCS: viper.GetBool("new.no-vcs"),
		Path:  viper.GetString("new.path"),
	}

	if err := service_repository.New(cfg, options); err != nil {
		return err
	}

	ui.Message(cfg, "New service repository", "✅ Project successfully created")
	return nil
}

func newProtobufModule(cfg *settings.Settings) error {
	options := &protobuf.NewOptions{
		Profile: viper.GetString("new.profile"),
	}

	return protobuf.New(cfg, options)
}

func newServiceTemplate(cfg *settings.Settings) error {
	options := &service.NewOptions{
		Path:          viper.GetString("new.path"),
		ProtoFilename: viper.GetString("new.proto"),
	}

	if err := service.New(cfg, options); err != nil {
		return err
	}

	ui.Message(cfg, "New service", "✅ Project successfully created")
	return nil
}

func setNewCmdFlags(cmd *cobra.Command) {
	// path option
	cmd.Flags().String("path", "", "Sets the output path name (default cwd).")
	_ = viper.BindPFlag("new.path", cmd.Flags().Lookup("path"))

	// proto file option
	cmd.Flags().String("proto", "", "Uses an _api.proto file as source for the service API.")
	_ = viper.BindPFlag("new.proto", cmd.Flags().Lookup("proto"))

	// no-vcs option
	cmd.Flags().Bool("no-vcs", false, "Disables creating projects with VCS support (default true).")
	_ = viper.BindPFlag("new.no-vcs", cmd.Flags().Lookup("no-vcs"))

	// profile option
	cmd.Flags().String("profile", "default", "Sets the profile to use.")
	_ = viper.BindPFlag("new.profile", cmd.Flags().Lookup("profile"))
}

func runNewProjectForm(cfg *settings.Settings) (string, error) {
	var selectedProject string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a project to create or Quit to exit the application").
				Options(
					huh.NewOption("Application/Service", "service-template"),
					huh.NewOption("Protobuf module", "protobuf-module"),
					huh.NewOption("Protobuf repository", "protobuf-monorepo"),
					huh.NewOption("Services repository", "services-monorepo"),
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
