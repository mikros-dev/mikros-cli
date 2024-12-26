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

	kind               Kind
	featureDefinitions map[string]*surveyAnswersDefinitions
	serviceDefinitions *surveyAnswersDefinitions
}

func newInitSurveyAnswers(kind Kind) *initSurveyAnswers {
	return &initSurveyAnswers{
		kind: kind,
	}
}

func (i *initSurveyAnswers) TemplateNames() []templates.TemplateFile {
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

func (i *initSurveyAnswers) AddFeatureDefinitions(name string, answers interface{}, save bool) {
	if i.featureDefinitions == nil {
		i.featureDefinitions = make(map[string]*surveyAnswersDefinitions)
	}

	i.featureDefinitions[name] = &surveyAnswersDefinitions{
		save:        save,
		definitions: answers,
	}
}

func (i *initSurveyAnswers) SetServiceDefinitions(answers interface{}, save bool) {
	i.serviceDefinitions = &surveyAnswersDefinitions{
		save:        save,
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
	save        bool
	definitions interface{}
}

func (s *surveyAnswersDefinitions) ShouldBeSaved() bool {
	return s.save
}

func (s *surveyAnswersDefinitions) Definitions() interface{} {
	return s.definitions
}
