package protobuf

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func New(cfg *settings.Settings) error {
	name, kind, err := chooseService(cfg)
	if err != nil {
		return err
	}

	switch kind {
	case "grpc":
		_, _, _, err := runGrpcForm(cfg)
		if err != nil {
			return err
		}

	case "http":
		if err := runHttpForm(cfg); err != nil {
			return err
		}
	}

	fmt.Println(name)
	return nil
}

func chooseService(cfg *settings.Settings) (string, string, error) {
	var (
		serviceName string
		serviceKind string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Service name. Enter the service name").
				Value(&serviceName).
				Validate(ui.IsEmpty("service name cannot be empty")),

			huh.NewSelect[string]().
				Title("Select the service type").
				Options(
					huh.NewOption("grpc", "grpc"),
					huh.NewOption("http", "http"),
				).
				Value(&serviceKind),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return "", "", err
	}

	return serviceName, serviceKind, nil
}

func runGrpcForm(cfg *settings.Settings) (string, bool, []string, error) {
	var (
		entityName  string
		defaultRPCs = true
		customRPCs  []string
		text        string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Entity name. Enter the service main entity name:").
				Validate(ui.IsEmpty("entity name cannot be empty")).
				Value(&entityName),

			huh.NewConfirm().
				Title("Use default CRUD RPCs for the service?").
				Value(&defaultRPCs),

			huh.NewText().
				Title("Enter the custom RPCs names (one per line)").
				Value(&text),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return "", false, nil, err
	}

	if text != "" {
		customRPCs = strings.Split(text, "\n")
	}

	return entityName, defaultRPCs, customRPCs, nil
}

func runHttpForm(cfg *settings.Settings) error {
	var (
		isAuthenticated bool
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Is the service authenticated?").
				Value(&isAuthenticated),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return err
	}

	return nil
}
