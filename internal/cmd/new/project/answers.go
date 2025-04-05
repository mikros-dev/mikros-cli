package project

import (
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

type surveyAnswers struct {
	RepositoryName string `survey:"repository_name"`
	ProjectName    string `survey:"project_name"`
	VcsPath        string `survey:"vcs_path"`
}

func newSurveyAnswers(cfg *settings.Settings, profileName string) *surveyAnswers {
	profile := surveyProfile(cfg, profileName)

	return &surveyAnswers{
		RepositoryName: profile.Project.ProtobufMonorepo.RepositoryName,
		ProjectName:    profile.Project.ProtobufMonorepo.ProjectName,
		VcsPath:        profile.Project.ProtobufMonorepo.VcsPath,
	}
}

func surveyProfile(cfg *settings.Settings, profile string) *settings.Profile {
	defaultValue := &cfg.App
	if profile == "default" {
		return defaultValue
	}

	d, ok := cfg.Profile[profile]
	if !ok {
		return defaultValue
	}

	return &d
}
