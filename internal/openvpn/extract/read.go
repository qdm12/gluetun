package extract

import (
	"io"
	"os"
	"strings"
)

func readCustomConfigLines(filepath string) (
	lines []string, err error) {
	file, err := os.Open(filepath)
	if err != nil {
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

	return strings.Split(string(b), "\n"), nil
}
