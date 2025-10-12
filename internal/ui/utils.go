package ui

import (
	"fmt"

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

// Message displays a message with the provided title and text.
func Message(cfg *settings.Settings, title, text string) {
	n := huh.NewNote().
		Title(title).
		Description(text + "\n").
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	_ = n.Init()
	fmt.Println(n.View())
}
