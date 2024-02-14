package env

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type WarnerMock struct {
	WarnFunc func(string)
}

func (wf WarnerMock) Warn(s string) {
	wf.WarnFunc(s)
}

func Test_readSecretFileAsStringPtr(t *testing.T) {
	t.Parallel()

	// Define test cases
	testCases := map[string]struct {
		setupEnv              func(tempDir string) []string // Function to setup environment variables
		secretPathEnvKey      string
		defaultSecretFileName string
		deprecatedEnvVars     []string
		expectedStringPtr     *string  // Expected result
		expectedUnsetEnvVars  []string // Expected environment variables to be unset
		expectError           bool     // Whether an error is expected
	}{
		"env_raw_secret": {
			setupEnv: func(_ string) []string {
				os.Setenv("SECRET", "mysecret")
				return []string{"SECRET"}
			},
			secretPathEnvKey:     "SECRET",
			expectedStringPtr:    ptrTo("mysecret"),
			expectedUnsetEnvVars: []string{"SECRET"},
		},
		"env_fallback_secret": {
			setupEnv: func(tempDir string) []string {
				secretFilePath := filepath.Join(tempDir, "secret_file")
				assert.NoError(t, os.WriteFile(secretFilePath, []byte("fallbacksecret"), os.ModePerm))
				return []string{}
			},
			secretPathEnvKey:      "SECRET",
			defaultSecretFileName: "secret_file",
			expectedStringPtr:     ptrTo("fallbacksecret"),
		},
		"env_variable_file_secret": {
			setupEnv: func(tempDir string) []string {
				secretFilePath := filepath.Join(tempDir, "secret_file")
				os.WriteFile(secretFilePath, []byte("filesecret"), os.ModePerm)
				os.Setenv("SECRET_FILE", secretFilePath)
				return []string{"SECRET_FILE"}
			},
			secretPathEnvKey:     "SECRET",
			expectedStringPtr:    ptrTo("filesecret"),
			expectedUnsetEnvVars: []string{"SECRET_FILE"},
		},
		"deprecated_env_variable": {
			setupEnv: func(tempDir string) []string {
				os.Setenv("DEPRECATED_SECRET", "deprecated")
				return []string{"DEPRECATED_SECRET"}
			},
			secretPathEnvKey:     "SECRET",
			deprecatedEnvVars:    []string{"DEPRECATED_SECRET"},
			expectedStringPtr:    ptrTo("deprecated"),
			expectedUnsetEnvVars: []string{"DEPRECATED_SECRET"},
		},
		"deprecated_env_variable_file": {
			setupEnv: func(tempDir string) []string {
				secretFilePath := filepath.Join(tempDir, "secret_file")
				os.WriteFile(secretFilePath, []byte("deprecatedfilesecret"), os.ModePerm)
				os.Setenv("DEPRECATED_SECRETFILE", secretFilePath)
				return []string{"DEPRECATED_SECRETFILE"}
			},
			secretPathEnvKey:     "SECRET",
			deprecatedEnvVars:    []string{"DEPRECATED_SECRETFILE"},
			expectedStringPtr:    ptrTo("deprecatedfilesecret"),
			expectedUnsetEnvVars: []string{"DEPRECATED_SECRETFILE"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			// t.Parallel() Cannot be used here because of the environment variables
			tempDir := t.TempDir()

			t.Logf("Running test case: %s", name)

			vars := []string{}
			// Setup environment variables
			if testCase.setupEnv != nil {
				vars = testCase.setupEnv(tempDir)
			}

			warner := WarnerMock{
				WarnFunc: func(message string) {
					t.Logf("WARN: %s", message)
				},
			}

			defaultSecretPath := filepath.Join(tempDir, testCase.defaultSecretFileName)
			source := New(warner) // Assuming Source is the type that implements readSecretFileAsStringPtr

			stringPtr, err := source.readSecretFileAsStringPtr(testCase.secretPathEnvKey, defaultSecretPath, testCase.deprecatedEnvVars)

			assert.Equal(t, testCase.expectedStringPtr, stringPtr)
			if testCase.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check if the specified environment variables are unset
			for _, varName := range testCase.expectedUnsetEnvVars {
				_, exists := os.LookupEnv(varName)
				assert.False(t, exists, "Expected %s to be unset", varName)
			}

			// Unset leftover environment variables
			// So other tests don't get affected
			for _, varName := range vars {
				os.Unsetenv(varName)
			}
		})
	}
}
