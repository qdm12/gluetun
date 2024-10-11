package secrets

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Source_Get(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		makeSource func(tempDir string) (source *Source, err error)
		key        string
		value      string
		isSet      bool
	}{
		"empty_key": {
			makeSource: func(tempDir string) (source *Source, err error) {
				return &Source{
					rootDirectory: tempDir,
					environ:       map[string]string{},
				}, nil
			},
		},
		"no_secret_file": {
			makeSource: func(tempDir string) (source *Source, err error) {
				return &Source{
					rootDirectory: tempDir,
					environ:       map[string]string{},
				}, nil
			},
			key: "test_file",
		},
		"empty_secret_file": {
			makeSource: func(tempDir string) (source *Source, err error) {
				secretFilepath := filepath.Join(tempDir, "test_file")
				const permission = fs.FileMode(0o600)
				err = os.WriteFile(secretFilepath, nil, permission)
				if err != nil {
					return nil, err
				}
				return &Source{
					rootDirectory: tempDir,
					environ:       map[string]string{},
				}, nil
			},
			key:   "test_file",
			isSet: true,
		},
		"default_secret_file": {
			makeSource: func(tempDir string) (source *Source, err error) {
				secretFilepath := filepath.Join(tempDir, "test_file")
				const permission = fs.FileMode(0o600)
				err = os.WriteFile(secretFilepath, []byte{'A'}, permission)
				if err != nil {
					return nil, err
				}
				return &Source{
					rootDirectory: tempDir,
					environ:       map[string]string{},
				}, nil
			},
			key:   "test_file",
			value: "A",
			isSet: true,
		},
		"env_specified_secret_file": {
			makeSource: func(tempDir string) (source *Source, err error) {
				secretFilepath := filepath.Join(tempDir, "test_file_custom")
				const permission = fs.FileMode(0o600)
				err = os.WriteFile(secretFilepath, []byte{'A'}, permission)
				if err != nil {
					return nil, err
				}
				return &Source{
					rootDirectory: tempDir,
					environ: map[string]string{
						"TEST_FILE_SECRETFILE": secretFilepath,
					},
				}, nil
			},
			key:   "test_file",
			value: "A",
			isSet: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			source, err := testCase.makeSource(t.TempDir())
			require.NoError(t, err)

			value, isSet := source.Get(testCase.key)
			assert.Equal(t, testCase.value, value)
			assert.Equal(t, testCase.isSet, isSet)
		})
	}
}
