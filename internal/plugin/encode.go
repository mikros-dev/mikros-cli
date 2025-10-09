package plugin

import (
	"github.com/mikros-dev/mikros-cli/internal/plugin/data"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

// Encoder is the mechanism that a plugin must use to return its data.
type Encoder struct {
	*data.PluginData
}

// NewEncoder creates a new Encoder instance.
func NewEncoder() *Encoder {
	return &Encoder{
		PluginData: &data.PluginData{},
	}
}

// SetName sets the name of the plugin.
func (e *Encoder) SetName(name string) {
	e.PluginData.Name = name
}

// SetUIName sets the UI plugin name.
func (e *Encoder) SetUIName(uiName string) {
	e.PluginData.UIName = uiName
}

// SetSurvey sets the plugin survey.
func (e *Encoder) SetSurvey(s *survey.Survey) {
	e.Survey = s
}

// SetAnswers sets the answers of the plugin.
func (e *Encoder) SetAnswers(answers map[string]interface{}) {
	e.Answers = answers
}

// SetTemplate sets the plugin template.
func (e *Encoder) SetTemplate(template *template.Template) {
	e.Template = template
}

// SetKind sets the plugin kind.
func (e *Encoder) SetKind(kind string) {
	e.Kind = kind
}

// SetError sets the error of the plugin.
func (e *Encoder) SetError(err error) {
	e.Error = err.Error()
}

// Output returns the plugin data.
func (e *Encoder) Output() error {
	return e.PluginData.Output()
}
