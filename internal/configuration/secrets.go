package configuration

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qdm12/golibs/params"
)

var (
	ErrGetSecretFilepath = errors.New("cannot get secret file path from env")
	ErrReadSecretFile    = errors.New("cannot read secret file")
	ErrSecretFileIsEmpty = errors.New("secret file is empty")
	ErrReadNonSecretFile = errors.New("cannot read non secret file")
	ErrFilesDoNotExist   = errors.New("files do not exist")
)

func cleanSuffix(value string) string {
	value = strings.TrimSuffix(value, "\n")
	value = strings.TrimSuffix(value, "\r")
	return value
}

func (r *reader) getFromEnvOrSecretFile(envKey string, compulsory bool, retroKeys []string) (value string, err error) {
	envOptions := []params.OptionSetter{
		params.Compulsory(), // to fallback on file reading
		params.CaseSensitiveValue(),
		params.Unset(),
		params.RetroKeys(retroKeys, r.onRetroActive),
	}
	value, envErr := r.env.Get(envKey, envOptions...)
	if envErr == nil {
		value = cleanSuffix(value)
		return value, nil
	}

	secretFilepathEnvKey := envKey + "_SECRETFILE"
	defaultSecretFile := "/run/secrets/" + strings.ToLower(envKey)
	filepath, err := r.env.Get(secretFilepathEnvKey,
		params.CaseSensitiveValue(),
		params.Default(defaultSecretFile),
	)
	if err != nil {
		return "", fmt.Errorf("%w: environment variable %s: %s",
			ErrGetSecretFilepath, secretFilepathEnvKey, err)
	}

	file, fileErr := os.OpenFile(filepath, os.O_RDONLY, 0)
	if os.IsNotExist(fileErr) {
		if compulsory {
			return "", fmt.Errorf("environment variable %s: %w", envKey, envErr)
		}
		return "", nil
	} else if fileErr != nil {
		return "", fmt.Errorf("%w: %s: %s", ErrReadSecretFile, filepath, fileErr)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %s", ErrReadSecretFile, filepath, err)
	}

	value = string(b)
	value = cleanSuffix(value)
	if compulsory && value == "" {
		return "", fmt.Errorf("%s: %w", filepath, ErrSecretFileIsEmpty)
	}

	return value, nil
}

// Tries to read from the secret file then the non secret file.
func (r *reader) getFromFileOrSecretFile(secretName, filepath string) (
	b []byte, err error) {
	defaultSecretFile := "/run/secrets/" + strings.ToLower(secretName)
	key := strings.ToUpper(secretName) + "_SECRETFILE"
	secretFilepath, err := r.env.Get(key,
		params.CaseSensitiveValue(),
		params.Default(defaultSecretFile),
	)
	if err != nil {
		return b, fmt.Errorf("environment variable %s: %w: %s", key, ErrGetSecretFilepath, err)
	}

	b, err = readFromFile(secretFilepath)
	if err != nil && !os.IsNotExist(err) {
		return b, fmt.Errorf("%w: %s", ErrReadSecretFile, err)
	} else if err == nil {
		return b, nil
	}

	// Secret file does not exist, try the non secret file
	b, err = readFromFile(filepath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%w: %s", ErrReadSecretFile, err)
	} else if err == nil {
		return b, nil
	}
	return nil, fmt.Errorf("%w: %s and %s", ErrFilesDoNotExist, secretFilepath, filepath)
}

func readFromFile(filepath string) (b []byte, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	b, err = io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	if err := file.Close(); err != nil {
		return nil, err
	}
	return b, nil
}
