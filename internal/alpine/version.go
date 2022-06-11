package alpine

import (
	"context"
	"io"
	"os"
	"strings"
)

func (a *Alpine) Version(ctx context.Context) (version string, err error) {
	file, err := os.OpenFile(a.alpineReleasePath, os.O_RDONLY, 0)
	if err != nil {
		return "", err
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	version = strings.ReplaceAll(string(b), "\n", "")
	return version, nil
}
