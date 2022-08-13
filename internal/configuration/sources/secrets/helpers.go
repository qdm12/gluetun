package secrets

import (
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/sources/files"
)

// getCleanedEnv returns an environment variable value with
// surrounding spaces and trailing new line characters removed.
func getCleanedEnv(envKey string) (value string) {
	value = os.Getenv(envKey)
	value = strings.TrimSpace(value)
	value = strings.TrimSuffix(value, "\r\n")
	value = strings.TrimSuffix(value, "\n")
	return value
}

func readSecretFileAsStringPtr(secretPathEnvKey, defaultSecretPath string) (
	stringPtr *string, err error) {
	path := getCleanedEnv(secretPathEnvKey)
	if path == "" {
		path = defaultSecretPath
	}
	return files.ReadFromFile(path)
}
