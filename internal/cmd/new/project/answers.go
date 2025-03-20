package project

import (
	"github.com/creasty/defaults"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

type surveyAnswers struct {
	RepositoryName string `survey:"repository_name" default:"protobuf-workspace"`
	ProjectName    string `survey:"project_name" default:"services"`
	VcsPath        string `survey:"vcs_path"`
}

func newSurveyAnswers(cfg *settings.Settings) *surveyAnswers {
	a := &surveyAnswers{}
	if err := defaults.Set(a); err != nil {
		// Without default values
		return a
	}

	a.VcsPath = cfg.Project.Template.VcsPath
	return a
}
