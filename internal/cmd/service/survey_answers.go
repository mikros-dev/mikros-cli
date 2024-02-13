package service

import (
	"github.com/somatech1/mikros-cli/internal/templates"
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

func (i initSurveyAnswers) TemplateNames() []templates.TemplateName {
	names := []templates.TemplateName{
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
		names = append(names, templates.TemplateName{
			Name:      "lifecycle",
			Extension: "go",
		})
	}

	return names
}
