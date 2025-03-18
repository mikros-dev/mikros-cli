package definitions

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"github.com/somatech1/mikros/components/definition"
)

type WriteOptions struct {
	NoValidation bool
}

// Write writes a mikros 'service.toml' file locally.
func Write(path string, defs *definition.Definitions, options ...*WriteOptions) error {
	if defs == nil {
		return errors.New("cannot handle nil options")
	}

	var opt *WriteOptions
	if len(options) != 0 {
		opt = options[0]
	}

	if opt != nil && !opt.NoValidation {
		if err := defs.Validate(); err != nil {
			return err
		}
	}

	file, err := os.Create(filepath.Join(path, "service.toml"))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	en := toml.NewEncoder(file)
	if err := en.Encode(defs); err != nil {
		return err
	}

	return nil
}

// AppendService appends a new section inside the 'service.toml' file to be
// loaded as settings for a specific service type.
func AppendService(path, serviceType string, serviceDefs interface{}) error {
	if serviceDefs == nil {
		return errors.New("cannot handle nil definitions")
	}

	filename := filepath.Join(path, "service.toml")
	if err := writeServiceDefinitions(filename, serviceType, serviceDefs); err != nil {
		return err
	}

	return nil
}

func writeServiceDefinitions(filename, name string, defs interface{}) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	line := fmt.Sprintf("\n[services.%v]\n", name)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	en := toml.NewEncoder(file)
	if err := en.Encode(defs); err != nil {
		return err
	}

	return nil
}

// AppendFeature appends a new section inside the 'service.toml' file to be
// loaded as settings for a specific feature.
func AppendFeature(path, featureName string, featureDefs interface{}) error {
	if featureDefs == nil {
		return errors.New("cannot handle nil definitions")
	}

	filename := filepath.Join(path, "service.toml")
	defs, err := loadCurrentFile(filename)
	if err != nil {
		return err
	}

	newFeatureDefs, err := featureDefsToMap(featureDefs)
	if err != nil {
		return err
	}

	features, ok := defs["features"]
	if !ok {
		defs["features"] = map[string]interface{}{
			featureName: newFeatureDefs,
		}
	}
	if ok {
		features := features.(map[string]interface{})
		features[featureName] = newFeatureDefs
		defs["features"] = features
	}

	if err := appendFeatureDefinitions(filename, defs); err != nil {
		return err
	}
	return nil
}

func loadCurrentFile(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if _, err := toml.DecodeFile(path, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func featureDefsToMap(featureDefs interface{}) (map[string]interface{}, error) {
	var b bytes.Buffer

	en := toml.NewEncoder(&b)
	if err := en.Encode(featureDefs); err != nil {
		return nil, err
	}

	defs := make(map[string]interface{})
	if err := toml.Unmarshal(b.Bytes(), &defs); err != nil {
		return nil, err
	}

	return defs, nil
}

func appendFeatureDefinitions(filename string, defs interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	en := toml.NewEncoder(file)
	if err := en.Encode(defs); err != nil {
		return err
	}

	return nil
}
