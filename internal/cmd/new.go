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

var (
	newCmd = &cobra.Command{
		Use:   "new",
		Short: "Create a new mikros project",
		Long:  "new helps creating a new mikros project",
	}
)

func newCmdInit(cfg *settings.Settings) {
	setNewCmdFlags()
	newCmd.Run = func(cmd *cobra.Command, args []string) {
		selected, err := runNewProjectForm(cfg)
		if err != nil {
			fmt.Println("new:", err)
			return
		}

		switch selected {
		case "protobuf-monorepo":
			newProtobufRepository(cfg)

		case "services-monorepo":
			newServiceRepository(cfg)

		case "protobuf-module":
			newProtobufModule(cfg)

		case "service-template":
			newServiceTemplate(cfg)

		case "quit":
			// Just quits
			return
		}
	}

	rootCmd.AddCommand(newCmd)
}

func newProtobufRepository(cfg *settings.Settings) {
	options := &protobuf_repository.NewOptions{
		NoVCS:   viper.GetBool("project-no-vcs"),
		Path:    viper.GetString("project-path"),
		Profile: viper.GetString("project-profile"),
	}

	if err := protobuf_repository.New(cfg, options); err != nil {
		fmt.Println("new:", err)
		return
	}

	fmt.Printf("\n✅ Project successfully created\n\n")
	fmt.Println("In order to start, execute the following command inside the new project directory:")
	fmt.Printf("\n$ make setup\n\n")
}

func newServiceRepository(cfg *settings.Settings) {
	options := &service_repository.NewOptions{
		NoVCS: viper.GetBool("project-no-vcs"),
		Path:  viper.GetString("project-path"),
	}

	if err := service_repository.New(cfg, options); err != nil {
		fmt.Println("new:", err)
		return
	}

	fmt.Printf("\n✅ Project successfully created\n\n")
}

func newProtobufModule(cfg *settings.Settings) {
	options := &protobuf.NewOptions{
		Profile: viper.GetString("project-profile"),
	}

	if err := protobuf.New(cfg, options); err != nil {
		fmt.Println("new:", err)
		return
	}
}

func newServiceTemplate(cfg *settings.Settings) {
	options := &service.NewOptions{
		Path:          viper.GetString("project-path"),
		ProtoFilename: viper.GetString("project-proto"),
	}

	if err := service.New(cfg, options); err != nil {
		fmt.Println("new:", err)
		return
	}

	fmt.Printf("\n✅ Service successfully created\n")
}

func setNewCmdFlags() {
	// path option
	newCmd.Flags().String("path", "", "Sets the output path name (default cwd).")
	_ = viper.BindPFlag("project-path", newCmd.Flags().Lookup("path"))

	// proto file option
	newCmd.Flags().String("proto", "", "Uses an _api.proto file as source for the service API.")
	_ = viper.BindPFlag("project-proto", newCmd.Flags().Lookup("proto"))

	// no-vcs option
	newCmd.Flags().Bool("no-vcs", false, "Disables creating projects with VCS support (default true).")
	_ = viper.BindPFlag("project-no-vcs", newCmd.Flags().Lookup("no-vcs"))

	// profile option
	newCmd.Flags().String("profile", "default", "Sets the profile to use.")
	_ = viper.BindPFlag("project-profile", newCmd.Flags().Lookup("profile"))
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
