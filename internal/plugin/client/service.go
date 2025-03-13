package client

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"

	"github.com/mikros-dev/mikros-cli/internal/plugin/data"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

type Service struct {
	name string
}

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
		return "", err
	}

	return out.String(), nil
}

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

func (s *Service) GetName() (string, error) {
	out, err := s.exec("-n")
	if err != nil {
		return "", err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return "", err
	}

	return d.Name, nil
}

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

func (s *Service) ValidateAnswers(answers map[string]interface{}) (map[string]interface{}, bool, error) {
	b, err := json.Marshal(answers)
	if err != nil {
		return nil, false, err
	}

	out, err := s.exec("-v", "-i", string(b))
	if err != nil {
		return nil, false, err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return nil, false, err
	}

	return d.Answers.Answers, d.Answers.Write, nil
}

func (s *Service) GetTemplates() (*template.Template, error) {
	out, err := s.exec("-t")
	if err != nil {
		return nil, err
	}

	d, err := data.DecodePluginData(out)
	if err != nil {
		return nil, err
	}

	return d.Template, nil
}
