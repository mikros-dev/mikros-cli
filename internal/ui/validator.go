package ui

import (
	"errors"
)

func IsEmpty(message string) func(s string) error {
	return func(s string) error {
		if s == "" {
			return errors.New(message)
		}

		return nil
	}
}
