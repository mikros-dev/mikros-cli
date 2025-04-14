package service_repository

import (
	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func runSurvey(cfg *settings.Settings) (*surveyAnswers, error) {
	answers := &surveyAnswers{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Repository name. Enter the name of the repository to create:").
				Value(&answers.RepositoryName).
				Validate(ui.IsEmpty("repository name cannot be empty")),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	return answers, nil
}
