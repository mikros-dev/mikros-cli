package service

import (
	"github.com/creasty/defaults"

	"github.com/mikros-dev/mikros-cli/internal/template"
)

type surveyAnswers struct {
	Name      string
	Type      string
	Language  string
	Version   string `default:"v0.1.0"`
	Product   string
	Features  []string
	Lifecycle []string

	featureDefinitions map[string]*surveyAnswersDefinitions
	serviceDefinitions *surveyAnswersDefinitions
}

func newSurveyAnswers() *surveyAnswers {
	a := &surveyAnswers{}
	if err := defaults.Set(a); err != nil {
		// Without default values
		return a
	}

	return a
}

func (i *surveyAnswers) TemplateNames() []template.File {
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

func (i *surveyAnswers) AddFeatureDefinitions(name string, answers interface{}) {
	if i.featureDefinitions == nil {
		i.featureDefinitions = make(map[string]*surveyAnswersDefinitions)
	}

	i.featureDefinitions[name] = &surveyAnswersDefinitions{
		definitions: answers,
	}
}

func (i *surveyAnswers) SetServiceDefinitions(answers interface{}) {
	i.serviceDefinitions = &surveyAnswersDefinitions{
		definitions: answers,
	}
}

func (i *surveyAnswers) ServiceDefinitions() *surveyAnswersDefinitions {
	return i.serviceDefinitions
}

func (i *surveyAnswers) FeatureDefinitions() map[string]*surveyAnswersDefinitions {
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
