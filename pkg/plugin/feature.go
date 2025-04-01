package plugin

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

// FeatureApi is the API that a feature plugin must implement to be supported
// by mikros CLI.
type FeatureApi interface {
	// Name must return the feature name that is registered inside the mikros
	// framework.
	Name() string

	// UIName must return the feature name that will be displayed for the user
	// during the survey to create a new service.
	UIName() string

	// Survey should return a survey.Survey object defining which properties
	// the user must configure to use this feature.
	Survey() *survey.Survey

	// ValidateAnswers receives answers from the feature survey to be validated
	// inside. It should return the data that should be written into the
	// 'service.toml' file.
	ValidateAnswers(in map[string]interface{}) (map[string]interface{}, error)
}

// Feature is the feature plugin object that provides the channel that mikros
// CLI recognizes as a plugin.
type Feature struct {
	api FeatureApi
}

// NewFeature creates a Feature object by receiving an object which must
// implement the FeatureApi interface.
func NewFeature(api FeatureApi) (*Feature, error) {
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

	if err := encoder.Output(); err != nil {
		return err
	}

	return nil
}

func inputToMap(in string) (map[string]interface{}, error) {
	var (
		out  map[string]interface{}
		data = strings.ReplaceAll(in, "\\", "")
	)

	if err := json.Unmarshal([]byte(data), &out); err != nil {
		return nil, fmt.Errorf("%v: %w", data, err)
	}

	return out, nil
}
