package auth

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Read reads the toml file specified by the filepath given.
// If the file does not exist, it returns empty settings and no error.
func Read(filepath string) (settings Settings, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Settings{}, nil
		}
		return settings, fmt.Errorf("opening file: %w", err)
	}
	decoder := toml.NewDecoder(file)
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&settings)
	if err == nil {
		return settings, nil
	}

	strictErr := new(toml.StrictMissingError)
	ok := errors.As(err, &strictErr)
	if !ok {
		return settings, fmt.Errorf("toml decoding file: %w", err)
	}
	return settings, fmt.Errorf("toml decoding file: %w:\n%s",
		strictErr, strictErr.String())
}
