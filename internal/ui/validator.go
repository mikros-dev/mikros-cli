package ui

import (
	"errors"
)

// IsEmpty returns a validator function that checks if a string is empty and
// returns an error with the given message if true.
func IsEmpty(message string) func(s string) error {
	return func(s string) error {
		if s == "" {
			return errors.New(message)
		}

		return nil
	}
}
