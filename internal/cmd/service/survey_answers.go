package service

import (
	"github.com/somatech1/mikros-cli/pkg/templates"
)

type initSurveyAnswers struct {
	Name      string
	Type      string
	Language  string
	Version   string
	Product   string
	Features  []string
	Lifecycle []string
}

func (i initSurveyAnswers) TemplateNames() []templates.TemplateFile {
	names := []templates.TemplateFile{
		{
			Name:      "main",
			Extension: "go",
		},
		{
			Name:      "service",
			Extension: "go",
		},
	}

	if len(i.Lifecycle) > 0 {
		names = append(names, templates.TemplateFile{
			Name:      "lifecycle",
			Extension: "go",
		})
	}

	return names
}
