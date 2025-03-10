package answers

import (
	"github.com/somatech1/mikros-cli/internal/templates"
)

type InitSurveyAnswers struct {
	Name          string
	Type          string
	Version       string
	Product       string
	Features      []string
	Lifecycle     []string
	RustLifecycle bool `survey:"lifecycle_rust"`

	language           templates.Language
	featureDefinitions map[string]*SurveyAnswersDefinitions
	serviceDefinitions *SurveyAnswersDefinitions
}

func NewInitSurveyAnswers(language templates.Language) *InitSurveyAnswers {
	return &InitSurveyAnswers{
		language: language,
	}
}

func (i *InitSurveyAnswers) Language() string {
	return i.language.String()
}

func (i *InitSurveyAnswers) AddFeatureDefinitions(name string, answers interface{}, save bool) {
	if i.featureDefinitions == nil {
		i.featureDefinitions = make(map[string]*SurveyAnswersDefinitions)
	}

	i.featureDefinitions[name] = &SurveyAnswersDefinitions{
		save:        save,
		definitions: answers,
	}
}

func (i *InitSurveyAnswers) SetServiceDefinitions(answers interface{}, save bool) {
	i.serviceDefinitions = &SurveyAnswersDefinitions{
		save:        save,
		definitions: answers,
	}
}

func (i *InitSurveyAnswers) ServiceDefinitions() *SurveyAnswersDefinitions {
	return i.serviceDefinitions
}

func (i *InitSurveyAnswers) FeatureDefinitions() map[string]*SurveyAnswersDefinitions {
	return i.featureDefinitions
}

type SurveyAnswersDefinitions struct {
	save        bool
	definitions interface{}
}

func (s *SurveyAnswersDefinitions) ShouldBeSaved() bool {
	return s.save
}

func (s *SurveyAnswersDefinitions) Definitions() interface{} {
	return s.definitions
}
