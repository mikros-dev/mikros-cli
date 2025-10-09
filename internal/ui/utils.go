package ui

import (
	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

// Alert displays a confirmation dialog with the provided text.
func Alert(cfg *settings.Settings, text string) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(text).
				Negative("").
				Affirmative("Ok"),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	return form.Run()
}
