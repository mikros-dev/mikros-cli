package protobuf

import (
	"fmt"

	"github.com/mikros-dev/mikros-cli/internal/settings"
)

func New(cfg *settings.Settings) error {
	name, kind, err := chooseService(cfg)
	if err != nil {
		return err
	}

	switch kind {
	case "grpc":
		_, _, _, err := runGrpcForm(cfg)
		if err != nil {
			return err
		}

	case "http":
		_, _, err := runHttpForm(cfg)
		if err != nil {
			return err
		}
	}

	fmt.Println(name)
	return nil
}
