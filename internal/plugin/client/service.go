package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/mikros-dev/mikros-cli/internal/plugin/data"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

// Service represents a service plugin.
type Service struct {
	name string
}

// NewService creates a new Service instance.
func NewService(path, name string) *Service {
	return &Service{
		name: filepath.Join(path, name),
	}
}

func (s *Service) exec(args ...string) (string, error) {
	cmd := exec.Command(s.name, args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		// Error here must be decoded from stdout
		d, err := data.DecodePluginData(out.String())
		if err != nil {
			// Nothing to do here, not our error
			return "", fmt.Errorf("error running service plugin: %w", err)
		}

		return "", errors.New(d.Error)
	}

	return out.String(), nil
}

// GetKind returns the kind of the service.
func (s *Service) GetKind() (string, error) {
	out, err := s.exec("-k")
	if err != nil {
		return "", err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return "", err
	}

	return d.Kind, nil
}

// GetSurvey returns the survey configuration associated with the service.
func (s *Service) GetSurvey() (*survey.Survey, error) {
	out, err := s.exec("-s")
	if err != nil {
		return nil, err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return nil, err
	}

	return d.Survey, nil
}

// ValidateAnswers validates the provided answers using the service plugin.
func (s *Service) ValidateAnswers(answers map[string]interface{}) (map[string]interface{}, error) {
	b, err := json.Marshal(answers)
	if err != nil {
		return nil, err
	}

	out, err := s.exec("-v", "-i", string(b))
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

// GetTemplates returns the templates associated with the service plugin.
func (s *Service) GetTemplates(answers map[string]interface{}) (*template.Template, error) {
	b, err := json.Marshal(answers)
	if err != nil {
		return nil, fmt.Errorf("error marshaling answers: %w", err)
	}

	out, err := s.exec("-t", "-i", string(b))
	if err != nil {
		return nil, err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return nil, err
	}

	return d.Template, nil
}
