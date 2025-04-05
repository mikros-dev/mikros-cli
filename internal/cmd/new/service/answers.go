package service

import (
	"github.com/creasty/defaults"
	"github.com/mikros-dev/mikros-cli/internal/protobuf"

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

	serviceAnswers     map[string]interface{}
	featureDefinitions map[string]*surveyAnswersDefinitions
	serviceDefinitions *surveyAnswersDefinitions
}

func newSurveyAnswers(protoFilename string) (*surveyAnswers, error) {
	a := &surveyAnswers{}
	if err := defaults.Set(a); err != nil {
		// Without default values
		return loadProtoValues(protoFilename, a)
	}

	return loadProtoValues(protoFilename, a)
}

func loadProtoValues(protoFilename string, a *surveyAnswers) (*surveyAnswers, error) {
	if protoFilename == "" {
		return a, nil
	}

	p, err := protobuf.Parse(protoFilename)
	if err != nil {
		return nil, err
	}

	a.Name = p.ServiceName
	return a, nil
}

func (s *surveyAnswers) TemplateNames() []template.File {
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

	if len(s.Lifecycle) > 0 {
		names = append(names, template.File{
			Name:      "lifecycle",
			Extension: "go",
		})
	}

	return names
}

func (s *surveyAnswers) AddFeatureDefinitions(name string, answers interface{}) {
	if s.featureDefinitions == nil {
		s.featureDefinitions = make(map[string]*surveyAnswersDefinitions)
	}

	s.featureDefinitions[name] = &surveyAnswersDefinitions{
		definitions: answers,
	}
}

func (s *surveyAnswers) SetServiceDefinitions(answers interface{}) {
	s.serviceDefinitions = &surveyAnswersDefinitions{
		definitions: answers,
	}
}

func (s *surveyAnswers) SetServiceAnswers(answers map[string]interface{}) {
	s.serviceAnswers = answers
}

func (s *surveyAnswers) ServiceDefinitions() *surveyAnswersDefinitions {
	return s.serviceDefinitions
}

func (s *surveyAnswers) FeatureDefinitions() map[string]*surveyAnswersDefinitions {
	return s.featureDefinitions
}

func (s *surveyAnswers) ServiceAnswers() map[string]interface{} {
	return s.serviceAnswers
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
