package edit

import (
	"fmt"

	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func runForm(cfg *settings.Settings) error {
loop:
	for {
		choice, err := initialForm(cfg)
		if err != nil {
			return err
		}

		switch choice {
		case "quit":
			break loop

		case "settings":
			if err := settingsForm(cfg); err != nil {
				return err
			}

		case "profiles":
			if err := profilesForm(cfg); err != nil {
				return err
			}
		}
	}

	return nil
}

func initialForm(cfg *settings.Settings) (string, error) {
	var option string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose the settings file section").
				Options(
					huh.NewOption("Adjust settings", "settings"),
					huh.NewOption("Profiles", "profiles"),
					huh.NewOption("Quit", "quit"),
				).
				Value(&option),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return option, nil
}

func settingsForm(cfg *settings.Settings) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Feature plugins:").Value(&cfg.Paths.Plugins.Features),
			huh.NewInput().Title("Service plugins:").Value(&cfg.Paths.Plugins.Services),
		).Title("Paths").Description("Configure paths for plugins\n"),

		huh.NewGroup(
			huh.NewConfirm().Title("Enable accessibility?").Value(&cfg.UI.Accessible),
			huh.NewSelect[string]().
				Title("Select the color theme to use:").
				Options(
					huh.NewOption("base16", "base16"),
					huh.NewOption("charm", "charm"),
					huh.NewOption("dracula", "dracula"),
					huh.NewOption("catppuccin", "catppuccin"),
					huh.NewOption("default", "default"),
				).
				Value(&cfg.UI.Theme),
		).Title("UI\n"),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return err
	}

	return nil
}

func profilesForm(cfg *settings.Settings) error {
loop:
	for {
		var (
			entries []huh.Option[string]
			choice  string
		)

		for k := range cfg.Profile {
			entries = append(entries, huh.NewOption(k, k))
		}
		entries = append(entries, []huh.Option[string]{
			huh.NewOption("Add new Profile", "add"),
			huh.NewOption("Remove Profile", "remove"),
			huh.NewOption("Back", "back"),
		}...)

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Choose the profile action").
					Options(entries...).
					Value(&choice),
			),
		).
			WithAccessible(cfg.UI.Accessible).
			WithTheme(cfg.GetTheme())

		if err := form.Run(); err != nil {
			return err
		}

		switch choice {
		case "back":
			break loop

		case "add":
			name, err := addProfile(cfg)
			if err != nil {
				return err
			}
			_, ok := cfg.Profile[name]
			if ok {
				if err := ui.Alert(cfg, fmt.Sprintf("profile '%s' already exists", name)); err != nil {
					return err
				}
			}

			if cfg.Profile == nil {
				cfg.Profile = make(map[string]settings.Profile)
			}
			cfg.Profile[name] = settings.Profile{}

		case "remove":
			if err := removeProfile(cfg); err != nil {
				return err
			}

		default:
			// edit current profile
			if err := editProfile(cfg, choice); err != nil {
				return err
			}
		}
	}

	return nil
}

func addProfile(cfg *settings.Settings) (string, error) {
	var name string
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Profile name. Enter the new profile name:").
				Validate(ui.IsEmpty("profile name cannot be empty")).
				Value(&name),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return "", err
	}

	return name, nil
}

func removeProfile(cfg *settings.Settings) error {
	var (
		names   []string
		entries []huh.Option[string]
	)

	for k := range cfg.Profile {
		entries = append(entries, huh.NewOption(k, k))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select profiles to remove").
				Options(entries...).
				Value(&names),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return err
	}
	if len(names) != 0 {
		for _, name := range names {
			delete(cfg.Profile, name)
		}
	}

	return nil
}

func editProfile(cfg *settings.Settings, name string) error {
	var (
		newName = name
		profile = cfg.Profile[name]
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Profile name. Enter the profile name:").
				Value(&newName).
				Validate(ui.IsEmpty("profile name cannot be empty")),

			huh.NewInput().
				Title("Repository name. Enter the name of the repository to create:").
				Value(&profile.Project.ProtobufMonorepo.RepositoryName).
				Validate(ui.IsEmpty("repository name cannot be empty")),

			huh.NewInput().
				Title("Project name. Enter your protobuf project name:").
				Value(&profile.Project.ProtobufMonorepo.ProjectName).
				Validate(ui.IsEmpty("project name cannot be empty")),

			huh.NewInput().
				Title("VCS path prefix. Enter your VCS path prefix to use for the project:").
				Value(&profile.Project.ProtobufMonorepo.VcsPath).
				Validate(ui.IsEmpty("VCS path prefix cannot be empty")),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if newName != name {
		delete(cfg.Profile, name)
	}
	cfg.Profile[newName] = profile

	return nil
}

func confirmSave(cfg *settings.Settings) error {
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Confirm saving the settings file?").
				Value(&confirm).
				Description("\nATTENTION! All settings will be overwritten."),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return err
	}

	if confirm {
		if err := cfg.Write(); err != nil {
			return err
		}
	}

	return nil
}
