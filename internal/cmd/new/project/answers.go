package project

import (
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

type surveyAnswers struct {
	RepositoryName string `survey:"repository_name"`
	ProjectName    string `survey:"project_name"`
	VcsPath        string `survey:"vcs_path"`
}

func newSurveyAnswers(cfg *settings.Settings, profile string) *surveyAnswers {
	values := surveyDefaultValues(cfg, profile)

	return &surveyAnswers{
		RepositoryName: values.RepositoryName,
		ProjectName:    values.ProjectName,
		VcsPath:        values.VcsPath,
	}
}

func surveyDefaultValues(cfg *settings.Settings, profile string) settings.ProtobufMonorepo {
	if profile == "default" {
		return cfg.Project.ProtobufMonorepo
	}

	d, ok := cfg.Profile[profile]
	if !ok {
		return cfg.Project.ProtobufMonorepo
	}

	return d.Project.ProtobufMonorepo
}
