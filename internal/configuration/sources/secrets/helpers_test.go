package secrets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/qdm12/gosettings/sources/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ptrTo[T any](value T) *T { return &value }

func Test_readSecretFileAsStringPtr(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		source                func(tempDir string) Source
		secretPathEnvKey      string
		defaultSecretFileName string
		setupFile             func(tempDir string) error
		stringPtr             *string
		errWrapped            error
		errMessage            string
	}{
		"no_secret_file": {
			defaultSecretFileName: "default_secret_file",
			secretPathEnvKey:      "SECRET_FILE",
		},
		"empty_secret_file": {
			defaultSecretFileName: "default_secret_file",
			secretPathEnvKey:      "SECRET_FILE",
			setupFile: func(tempDir string) error {
				secretFilepath := filepath.Join(tempDir, "default_secret_file")
				return os.WriteFile(secretFilepath, nil, os.ModePerm)
			},
			stringPtr: ptrTo(""),
		},
		"default_secret_file": {
			defaultSecretFileName: "default_secret_file",
			secretPathEnvKey:      "SECRET_FILE",
			setupFile: func(tempDir string) error {
				secretFilepath := filepath.Join(tempDir, "default_secret_file")
				return os.WriteFile(secretFilepath, []byte("A"), os.ModePerm)
			},
			stringPtr: ptrTo("A"),
		},
		"env_specified_secret_file": {
			source: func(tempDir string) Source {
				secretFilepath := filepath.Join(tempDir, "secret_file")
				environ := []string{"SECRET_FILE=" + secretFilepath}
				return Source{env: *env.New(environ)}
			},
			defaultSecretFileName: "default_secret_file",
			secretPathEnvKey:      "SECRET_FILE",
			setupFile: func(tempDir string) error {
				secretFilepath := filepath.Join(tempDir, "secret_file")
				return os.WriteFile(secretFilepath, []byte("B"), os.ModePerm)
			},
			stringPtr: ptrTo("B"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tempDir := t.TempDir()

			var source Source
			if testCase.source != nil {
				source = testCase.source(tempDir)
			}

			defaultSecretPath := filepath.Join(tempDir, testCase.defaultSecretFileName)
			if testCase.setupFile != nil {
				err := testCase.setupFile(tempDir)
				require.NoError(t, err)
			}

			stringPtr, err := source.readSecretFileAsStringPtr(
				testCase.secretPathEnvKey, defaultSecretPath)

			assert.Equal(t, testCase.stringPtr, stringPtr)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
