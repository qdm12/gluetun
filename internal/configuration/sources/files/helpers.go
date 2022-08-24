package files

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/openvpn/extract"
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
	content = strings.TrimSuffix(content, "\r\n")
	content = strings.TrimSuffix(content, "\n")
	return &content, nil
}

func readPEMFile(filepath string) (base64Ptr *string, err error) {
	pemData, err := ReadFromFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	if pemData == nil {
		return nil, nil //nolint:nilnil
	}

	base64Data, err := extract.PEM([]byte(*pemData))
	if err != nil {
		return nil, fmt.Errorf("extracting base64 encoded data from PEM content: %w", err)
	}

	return &base64Data, nil
}
