package plugin

import (
	"errors"
	"flag"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

type ServiceApi interface {
	Name() string
	Kind() string
	Survey() *survey.Survey
	ValidateAnswers(in map[string]interface{}) (map[string]interface{}, bool, error)
	Template() *template.Template
}

type Service struct {
	api ServiceApi
}

func NewService(api ServiceApi) (*Service, error) {
	if api == nil {
		return nil, errors.New("api cannot be nil")
	}

	return &Service{
		api: api,
	}, nil
}

func (s *Service) Run() error {
	// Supported plugin options
	nFlag := flag.Bool("n", false, "Get plugin name")
	sFlag := flag.Bool("s", false, "Retrieve feature survey questions")
	vFlag := flag.Bool("v", false, "Validate answers")
	tFlag := flag.Bool("t", false, "Retrieve plugin custom templates")
	kFlag := flag.Bool("k", false, "Get service kind")
	input := flag.String("i", "", "Input values for plugin arguments")
	flag.Parse()

	encoder := plugin.NewEncoder()

	switch {
	case *nFlag:
		encoder.SetName(s.api.Name())
	case *sFlag:
		encoder.SetSurvey(s.api.Survey())
	case *vFlag:
		if *input != "" {
			return errors.New("invalid input")
		}

		in, err := inputToMap(*input)
		if err != nil {
			return err
		}

		data, save, err := s.api.ValidateAnswers(in)
		if err != nil {
			return err
		}

		encoder.SetAnswers(data, save)
	case *tFlag:
		encoder.SetTemplate(s.api.Template())
	case *kFlag:
		encoder.SetKind(s.api.Kind())
	default:
		return errors.New("no valid command specified")
	}

	if err := encoder.Output(); err != nil {
		return err
	}

	return nil
}
