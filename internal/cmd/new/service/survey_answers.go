package service

import (
	"github.com/mikros-dev/mikros-cli/internal/template"
)

type initSurveyAnswers struct {
	Name      string
	Type      string
	Language  string
	Version   string
	Product   string
	Features  []string
	Lifecycle []string

	featureDefinitions map[string]*surveyAnswersDefinitions
	serviceDefinitions *surveyAnswersDefinitions
}

func (i *initSurveyAnswers) TemplateNames() []template.File {
	names := []template.File{
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
		names = append(names, template.File{
			Name:      "lifecycle",
			Extension: "go",
		})
	}

	return names
}

func (i *initSurveyAnswers) AddFeatureDefinitions(name string, answers interface{}) {
	if i.featureDefinitions == nil {
		i.featureDefinitions = make(map[string]*surveyAnswersDefinitions)
	}

	i.featureDefinitions[name] = &surveyAnswersDefinitions{
		definitions: answers,
	}
}

func (i *initSurveyAnswers) SetServiceDefinitions(answers interface{}) {
	i.serviceDefinitions = &surveyAnswersDefinitions{
		definitions: answers,
	}
}

func (i *initSurveyAnswers) ServiceDefinitions() *surveyAnswersDefinitions {
	return i.serviceDefinitions
}

func (i *initSurveyAnswers) FeatureDefinitions() map[string]*surveyAnswersDefinitions {
	return i.featureDefinitions
}

type surveyAnswersDefinitions struct {
	definitions interface{}
}

func (s *surveyAnswersDefinitions) ShouldBeSaved() bool {
	return s != nil && s.definitions != nil
}

func (s *surveyAnswersDefinitions) Definitions() interface{} {
	return s.definitions
}
