package plugin

import (
	"errors"
	"flag"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
	"github.com/mikros-dev/mikros-cli/pkg/template"
)

// ServiceApi is the API that a service plugin must implement to be supported
// by mikros CLI.
type ServiceApi interface {
	// Kind must return the new service type for services.
	Kind() string

	// Survey should return a survey.Survey object defining which properties
	// the user must configure to use this service type.
	Survey() *survey.Survey

	// ValidateAnswers receives answers from the service survey to be validated
	// inside. It should return the data that should be written into the
	// 'service.toml' file and a flag indicating if it should be written or not.
	ValidateAnswers(in map[string]interface{}) (map[string]interface{}, bool, error)
	Template() *template.Template
}

// Service is the service plugin object that provides the channel that mikros
// CLI recognizes as a plugin.
type Service struct {
	api ServiceApi
}

// NewService creates a Service object by receiving an object which must
// implement the ServiceApi interface.
func NewService(api ServiceApi) (*Service, error) {
	if api == nil {
		return nil, errors.New("api cannot be nil")
	}

	return &Service{
		api: api,
	}, nil
}

// Run executes the plugin.
func (s *Service) Run() error {
	// Supported plugin options
	sFlag := flag.Bool("s", false, "Retrieve feature survey questions")
	vFlag := flag.Bool("v", false, "Validate answers")
	tFlag := flag.Bool("t", false, "Retrieve plugin custom templates")
	kFlag := flag.Bool("k", false, "Get service kind")
	input := flag.String("i", "", "Input values for plugin arguments")
	flag.Parse()

	encoder := plugin.NewEncoder()

	switch {
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
