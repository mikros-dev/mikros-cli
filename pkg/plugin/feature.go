package plugin

import (
	"encoding/json"
	"errors"
	"flag"

	"github.com/mikros-dev/mikros-cli/internal/plugin"
	"github.com/mikros-dev/mikros-cli/pkg/survey"
)

type FeatureApi interface {
	Name() string
	UIName() string
	Survey() *survey.Survey
	ValidateAnswers(in map[string]interface{}) (map[string]interface{}, bool, error)
}

type Feature struct {
	api FeatureApi
}

func NewFeature(api FeatureApi) (*Feature, error) {
	if api == nil {
		return nil, errors.New("api cannot be nil")
	}

	return &Feature{
		api: api,
	}, nil
}

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
		if *input != "" {
			return errors.New("invalid input")
		}

		in, err := inputToMap(*input)
		if err != nil {
			return err
		}

		data, save, err := f.api.ValidateAnswers(in)
		if err != nil {
			return err
		}

		encoder.SetAnswers(data, save)
	default:
		return errors.New("no valid command specified")
	}

	if err := encoder.Output(); err != nil {
		return err
	}

	return nil
}

func inputToMap(in string) (map[string]interface{}, error) {
	var out map[string]interface{}

	if err := json.Unmarshal([]byte(in), &out); err != nil {
		return nil, err
	}

	return out, nil
}
