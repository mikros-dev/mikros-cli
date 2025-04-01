package plugin

import (
	"github.com/mikros-dev/mikros-cli/internal/plugin/data"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

type Encoder struct {
	*data.PluginData
}

func NewEncoder() *Encoder {
	return &Encoder{
		PluginData: &data.PluginData{},
	}
}

func (e *Encoder) SetName(name string) {
	e.PluginData.Name = name
}

func (e *Encoder) SetUIName(uiName string) {
	e.PluginData.UIName = uiName
}

func (e *Encoder) SetSurvey(s *survey.Survey) {
	e.Survey = s
}

func (e *Encoder) SetAnswers(answers map[string]interface{}) {
	e.Answers = answers
}

func (e *Encoder) SetTemplate(template *template.Template) {
	e.Template = template
}

func (e *Encoder) SetKind(kind string) {
	e.Kind = kind
}

func (e *Encoder) SetError(err error) {
	e.Error = err.Error()
}

func (e *Encoder) Output() error {
	return e.PluginData.Output()
}
