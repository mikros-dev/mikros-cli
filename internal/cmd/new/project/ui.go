package project

import (
	"github.com/charmbracelet/huh"

	"github.com/mikros-dev/mikros-cli/internal/settings"
	"github.com/mikros-dev/mikros-cli/internal/ui"
)

func runSurvey(cfg *settings.Settings, profile string) (*surveyAnswers, error) {
	answers := newSurveyAnswers(cfg, profile)
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Repository name. Enter the name of the repository to create:").
				Value(&answers.RepositoryName).
				Validate(ui.IsEmpty("repository name cannot be empty")),

			huh.NewInput().
				Title("Project name. Enter your protobuf project name:").
				Value(&answers.ProjectName).
				Validate(ui.IsEmpty("project name cannot be empty")),

			huh.NewInput().
				Title("VCS path prefix. Enter your VCS path prefix to use for the project:").
				Value(&answers.VcsPath).
				Validate(ui.IsEmpty("VCS path prefix cannot be empty")),
		),
	).
		WithAccessible(cfg.UI.Accessible).
		WithTheme(cfg.GetTheme())

	if err := form.Run(); err != nil {
		return nil, err
	}

	return answers, nil
}
