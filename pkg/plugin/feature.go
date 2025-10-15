package plugin

import (
	"errors"
	"flag"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
)

// FeatureAPI is the API that a feature plugin must implement to be supported
// by mikros CLI.
type FeatureAPI interface {
	// Name must return the feature name registered inside the mikros
	// framework.
	Name() string

	// UIName must return the feature name that will be displayed for the user
	// during the survey to create a new service.
	UIName() string

	// Survey should return a survey.Survey object defining which properties
	// the user must configure to use this feature.
	Survey() *Survey

	// ValidateAnswers receives answers from the feature survey to be validated
	// inside. It should return the data that should be written into the
	// 'service.toml' file.
	ValidateAnswers(in map[string]interface{}) (map[string]interface{}, error)
}

// Feature is the feature plugin object that provides the channel that mikros
// CLI recognizes as a plugin.
type Feature struct {
	api FeatureAPI
}

// NewFeature creates a Feature object by receiving an object which must
// implement the FeatureAPI interface.
func NewFeature(api FeatureAPI) (*Feature, error) {
	if api == nil {
		return nil, errors.New("api cannot be nil")
	}

	return &Feature{
		api: api,
	}, nil
}

// Run executes the plugin.
func (f *Feature) Run() error {
	// Supported plugin options
	nFlag := flag.Bool("n", false, "Get plugin name")
	uFlag := flag.Bool("u", false, "Get UI name")
	sFlag := flag.Bool("s", false, "Retrieve feature survey questions")
	vFlag := flag.Bool("v", false, "Validate answers")
	input := flag.String("i", "", "Input values for plugin arguments")
	flag.Parse()

	encoder := plugin.NewEncoder()

	switch {
	case *nFlag:
		encoder.SetName(f.api.Name())
	case *uFlag:
		encoder.SetUIName(f.api.UIName())
	case *sFlag:
		encoder.SetSurvey(f.api.Survey())
	case *vFlag:
		if *input == "" {
			return errors.New("invalid input")
		}

		in, err := inputToMap(*input)
		if err != nil {
			return err
		}

		data, err := f.api.ValidateAnswers(in)
		if err != nil {
			return err
		}

		encoder.SetAnswers(data)
	default:
		return errors.New("no valid command specified")
	}

	return encoder.Output()
}
