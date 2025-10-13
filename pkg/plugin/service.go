package plugin

import (
	"errors"
	"flag"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
)

// ServiceAPI is the API that a service plugin must implement to be supported
// by mikros CLI.
type ServiceAPI interface {
	// Kind must return the new service type for services.
	Kind() string

	// Survey should return a survey.Survey object defining which properties
	// the user must configure to use this service type.
	Survey() *Survey

	// ValidateAnswers receives answers from the service survey to be validated
	// inside. It should return the data that should be written into the
	// 'service.toml' file.
	ValidateAnswers(in map[string]interface{}) (map[string]interface{}, error)

	// Template allows the plugin to return a set of custom templates that will
	// be executed when a service is created. It also receives the answers from
	// the service survey.
	Template(in map[string]interface{}) *Template
}

// Service is the service plugin object that provides the channel that mikros
// CLI recognizes as a plugin.
type Service struct {
	api ServiceAPI
}

// NewService creates a Service object by receiving an object which must
// implement the ServiceAPI interface.
func NewService(api ServiceAPI) (*Service, error) {
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
		if *input == "" {
			return errors.New("invalid input")
		}
		in, err := inputToMap(*input)
		if err != nil {
			return err
		}

		data, err := s.api.ValidateAnswers(in)
		if err != nil {
			return err
		}

		encoder.SetAnswers(data)
	case *tFlag:
		if *input == "" {
			return errors.New("invalid input")
		}
		in, err := inputToMap(*input)
		if err != nil {
			return err
		}

		encoder.SetTemplate(s.api.Template(in))
	case *kFlag:
		encoder.SetKind(s.api.Kind())
	default:
		return errors.New("no valid command specified")
	}

	return encoder.Output()
}
