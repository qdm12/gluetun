package files

import (
	"io"
	"os"
)

// ReadFromFile reads the content of the file as a string.
// It returns a nil string pointer if the file does not exist.
func ReadFromFile(filepath string) (s *string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil //nolint:nilnil
		}
		return nil, err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	content := string(b)
	return &content, nil
}
