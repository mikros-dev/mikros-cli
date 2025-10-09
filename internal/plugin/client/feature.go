package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/mikros-dev/mikros-cli/internal/plugin/data"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

// Feature represents a feature plugin.
type Feature struct {
	name string
}

// NewFeature creates a new Feature instance.
func NewFeature(path, name string) *Feature {
	return &Feature{
		name: filepath.Join(path, name),
	}
}

func (f *Feature) exec(args ...string) (string, error) {
	cmd := exec.Command(f.name, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// Error here must be decoded from stdout
		d, err := data.DecodePluginData(out.String())
		if err != nil {
			// Nothing to do here, not our error
			return "", err
		}

		return "", errors.New(d.Error)
	}

	return out.String(), nil
}

// GetName retrieves the name of the plugin.
func (f *Feature) GetName() (string, error) {
	out, err := f.exec("-n")
	if err != nil {
		return "", err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return "", err
	}

	return d.Name, nil
}

// GetUIName retrieves the UI display name of the feature plugin.
func (f *Feature) GetUIName() (string, error) {
	out, err := f.exec("-u")
	if err != nil {
		return "", err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return "", err
	}

	return d.UIName, nil
}

// GetSurvey retrieves the survey configuration associated with the feature plugin.
func (f *Feature) GetSurvey() (*survey.Survey, error) {
	out, err := f.exec("-s")
	if err != nil {
		return nil, err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return nil, err
	}

	return d.Survey, nil
}

// ValidateAnswers validates the provided answers using the feature plugin.
func (f *Feature) ValidateAnswers(answers map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(answers)
	if err != nil {
		return nil, err
	}

	out, err := f.exec("-v", "-i", string(b))
	if err != nil {
		return nil, err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return nil, err
	}
	if len(d.Answers) == 0 {
		return nil, nil
	}

	return d.Answers, nil
}
