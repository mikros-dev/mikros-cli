package plugin

import (
	"encoding/json"

	"github.com/mikros-dev/mikros-cli/internal/plugin/wire"
)

// Encoder is the mechanism that a plugin must use to return its data.
type Encoder struct {
	*wire.PluginData
}

// NewEncoder creates a new Encoder instance.
func NewEncoder() *Encoder {
	return &Encoder{
		PluginData: &wire.PluginData{},
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
func (e *Encoder) SetSurvey(s interface{}) {
	if s == nil {
		return
	}

	b, err := json.Marshal(s)
	if err != nil {
		e.Error = err.Error()
		return
	}

	e.PluginData.Survey = b
}

// SetAnswers sets the answers of the plugin.
func (e *Encoder) SetAnswers(answers map[string]interface{}) {
	e.Answers = answers
}

// SetTemplate sets the plugin template.
func (e *Encoder) SetTemplate(t interface{}) {
	if t == nil {
		return
	}

	b, err := json.Marshal(t)
	if err != nil {
		e.Error = err.Error()
		return
	}

	e.PluginData.Template = b
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
