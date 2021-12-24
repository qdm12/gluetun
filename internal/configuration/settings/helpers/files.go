package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	ErrFileDoesNotExist = errors.New("file does not exist")
	ErrFileRead         = errors.New("cannot read file")
	ErrFileClose        = errors.New("cannot close file")
)

func FileExists(path string) (err error) {
	path = filepath.Clean(path)

	f, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %s", ErrFileDoesNotExist, path)
	} else if err != nil {
		return fmt.Errorf("%w: %s", ErrFileRead, err)
	}

	if err := f.Close(); err != nil {
		return fmt.Errorf("%w: %s", ErrFileClose, err)
	}

	return nil
}
