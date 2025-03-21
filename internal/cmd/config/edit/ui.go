package edit

import (
	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func runForm(cfg *settings.Settings) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Feature plugins:").Value(&cfg.Paths.Plugins.Features),
			huh.NewInput().Title("Service plugins:").Value(&cfg.Paths.Plugins.Services),
		).Title("Paths"),

		huh.NewGroup(
			huh.NewConfirm().Title("Enable accessibility?").Value(&cfg.UI.Accessible),
		).Title("UI"),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return err
	}

	return nil
}
