package env

import (
	"fmt"
	"os"
)

func unsetEnvKeys(envKeys []string, err error) (newErr error) {
	newErr = err
	for _, envKey := range envKeys {
		unsetErr := os.Unsetenv(envKey)
		if unsetErr != nil && newErr == nil {
			newErr = fmt.Errorf("unsetting environment variable %s: %w", envKey, unsetErr)
		}
	}
	return newErr
}

func ptrTo[T any](value T) *T {
	return &value
}
