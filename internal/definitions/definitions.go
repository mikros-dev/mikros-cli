package definitions

import (
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
func AppendService(path, name string, serviceDefs interface{}) error {
	if serviceDefs == nil {
		return errors.New("cannot handle nil definitions")
	}

	filename := filepath.Join(path, "service.toml")
	if err := writeServiceDefinitions(filename, name, serviceDefs); err != nil {
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

	line := fmt.Sprintf("\n[%v]\n", name)
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
func AppendFeature(path, name string, featureDefs interface{}) error {
	if featureDefs == nil {
		return errors.New("cannot handle nil definitions")
	}

	filename := filepath.Join(path, "service.toml")
	defs, err := loadCurrentFile(filename)
	if err != nil {
		return err
	}

	features, ok := defs["features"]
	if !ok {
		if err := writeFirstFeatureDefinitions(filename, name, featureDefs); err != nil {
			return err
		}
	}
	if ok {
		features := features.(map[string]interface{})
		features[name] = featureDefs
		defs["features"] = features

		if err := appendFeatureDefinitions(filename, defs); err != nil {
			return err
		}
	}

	return nil
}

func writeFirstFeatureDefinitions(filename, name string, defs interface{}) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	line := fmt.Sprintf("\n[features]\n[features.%v]\n", name)
	if _, err := file.WriteString(line); err != nil {
		return err
	}

	en := toml.NewEncoder(file)
	if err := en.Encode(defs); err != nil {
		return err
	}

	return nil
}

func appendFeatureDefinitions(filename string, defs map[string]interface{}) error {
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

func loadCurrentFile(path string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if _, err := toml.DecodeFile(path, &data); err != nil {
		return nil, err
	}

	return data, nil
}
