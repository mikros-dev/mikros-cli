package data

import (
	"encoding/json"
	"fmt"
	"github.com/mikros-dev/mikros-cli/pkg/template"

	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

type PluginData struct {
	Name     string             `json:"name,omitempty"`
	UIName   string             `json:"ui_name,omitempty"`
	Kind     string             `json:"kind,omitempty"`
	Survey   *survey.Survey     `json:"survey,omitempty"`
	Answers  *Answers           `json:"answers,omitempty"`
	Template *template.Template `json:"template,omitempty"`
	Error    error              `json:"error,omitempty"`
}

type Answers struct {
	Answers map[string]interface{} `json:"answers,omitempty"`
	Write   bool                   `json:"write,omitempty"`
}

func (p *PluginData) Output() error {
	b, err := json.Marshal(p)
	if err != nil {
		return err
	}

	fmt.Println(string(b))
	return nil
}

func DecodePluginData(in string) (*PluginData, error) {
	var p PluginData
	if err := json.Unmarshal([]byte(in), &p); err != nil {
		return nil, err
	}

	return &p, nil
}
