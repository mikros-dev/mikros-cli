package project

import (
	"github.com/mikros-dev/mikros-cli/internal/settings"
)

type surveyAnswers struct {
	RepositoryName string `survey:"repository_name"`
	ProjectName    string `survey:"project_name"`
	VcsPath        string `survey:"vcs_path"`
}

func newSurveyAnswers(cfg *settings.Settings) *surveyAnswers {
	a := &surveyAnswers{}

	// Use settings values as the default one
	a.VcsPath = cfg.Project.ProtobufMonorepo.VcsPath
	a.ProjectName = cfg.Project.ProtobufMonorepo.ProjectName
	a.RepositoryName = cfg.Project.ProtobufMonorepo.RepositoryName

	return a
}
