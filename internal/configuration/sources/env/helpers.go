package env

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
	"github.com/qdm12/gluetun/internal/openvpn/extract"
	"github.com/qdm12/gosettings/sources/env"
)

// readSecretFileAsStringPtr reads a secret provided by the user.
// It checks in the following order:
// 1. The environment variable key.
// 2. The envionment variable key with a suffix of _FILE.
// 3. The default file path.
// 4. The deprecated environment variables.
//
// At the end it unsets all the environment variables for safety.
//
// Parameters:
// - secretPathEnvKey: the environment variable key to check first.
// - defaultSecretPath: the default file path to check if the environment variable is not set.
// - deprecated_env_vars: the deprecated environment variables to check if the environment variable is not set.
//
// Returns:
// - A pointer to a string containing the file content, or nil if the file is empty or not specified.
// - An error if reading the file fails.
func (s *Source) readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath string, deprecatedEnvVars []string) (
	stringPtr *string, err error) {
	// No matter what happens, we want to unset all the environment variables
	defer func() {
		vars_to_unset := append(deprecatedEnvVars, secretPathEnvKey, secretPathEnvKey+"_FILE")
		err = unsetEnvKeys(vars_to_unset, err)
	}()

	// Try to get the secret from the environment variable.
	secretPathOrRaw := s.env.String(secretPathEnvKey, env.ForceLowercase(false))
	if secretPathOrRaw != "" {
		return &secretPathOrRaw, nil
	}

	// Try to get the secret from the environment variable with _FILE suffix.
	secretPathOrRaw = os.Getenv(secretPathEnvKey + "_FILE")
	if secretPathOrRaw != "" {
		return files.ReadFromFile(secretPathOrRaw)
	}

	// Check deprecated environment variables.
	for _, depVar := range deprecatedEnvVars {
		secretPathOrRaw = os.Getenv(depVar)
		if secretPathOrRaw != "" {
			s.warner.Warn(fmt.Sprintf("using deprecated environment variable %s, please use %s instead", depVar, secretPathEnvKey))

			// check if it ends with SECRETFILE
			if strings.HasSuffix(depVar, "SECRETFILE") {
				return files.ReadFromFile(secretPathOrRaw)
			}

			// raw deprecated secret
			return &secretPathOrRaw, nil
		}
	}

	// Use the default secret path if the environment variables are not set.
	if _, err := os.Stat(defaultSecretPath); err == nil {
		return files.ReadFromFile(defaultSecretPath)
	}

	return nil, errors.New("no secret file specified or found")
}

func (s *Source) readPEMSecretFile(secretPathEnvKey, defaultSecretPath string, deprecated_env_vars []string) (
	base64Ptr *string, err error) {
	pemData, err := s.readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath, deprecated_env_vars)
	if err != nil {
		return nil, fmt.Errorf("reading secret file: %w", err)
	}

	if pemData == nil {
		return nil, nil //nolint:nilnil
	}

	base64Data, err := extract.PEM([]byte(*pemData))
	if err != nil {
		return nil, fmt.Errorf("extracting base64 encoded data from PEM content: %w", err)
	}

	return &base64Data, nil
}

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

func firstKeySet(e env.Env, keys ...string) (firstKeySet string) {
	for _, key := range keys {
		value := e.Get(key)
		if value != nil {
			return key
		}
	}
	return ""
}
