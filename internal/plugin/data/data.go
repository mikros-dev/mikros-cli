package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

// PluginData represents the data that a plugin must return to the CLI.
type PluginData struct {
	Name     string                 `json:"name,omitempty"`
	UIName   string                 `json:"ui_name,omitempty"`
	Kind     string                 `json:"kind,omitempty"`
	Survey   *survey.Survey         `json:"survey,omitempty"`
	Answers  map[string]interface{} `json:"answers,omitempty"`
	Template *template.Template     `json:"template,omitempty"`
	Error    string                 `json:"error,omitempty"`
}

// Output prints the data to the CLI.
func (p *PluginData) Output() error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

// DecodePluginData decodes the data from the CLI.
func DecodePluginData(in string) (*PluginData, error) {
	var (
		d = json.NewDecoder(strings.NewReader(in))
		p PluginData
	)

	d.UseNumber()
	if err := d.Decode(&p); err != nil {
		return nil, err
	}

	return &p, nil
}
