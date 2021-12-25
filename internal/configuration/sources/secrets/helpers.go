package secrets

import (
	"os"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

func readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath string) (
	stringPtr *string, err error) {
	path := os.Getenv(secretPathEnvKey)
	if path == "" {
		path = defaultSecretPath
	}
	return files.ReadFromFile(path)
}

func readSecretFileAsString(secretPathEnvKey, defaultSecretPath string) (
	s string, err error) {
	path := os.Getenv(secretPathEnvKey)
	if path == "" {
		path = defaultSecretPath
	}
	stringPtr, err := files.ReadFromFile(path)
	if err != nil {
		return "", err
	} else if stringPtr == nil {
		return "", nil
	}
	return *stringPtr, nil
}
