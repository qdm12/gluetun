package files

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/qdm12/gluetun/internal/openvpn/extract"
)

// ReadFromFile reads the content of the file as a string,
// and returns if the file was present or not with isSet.
func ReadFromFile(filepath string) (content string, isSet bool, err error) {
	file, err := os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, fmt.Errorf("opening file: %w", err)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return "", false, fmt.Errorf("reading file: %w", err)
	}

	if err := file.Close(); err != nil {
		return "", false, fmt.Errorf("closing file: %w", err)
	}

	content = string(b)
	content = strings.TrimSuffix(content, "\r\n")
	content = strings.TrimSuffix(content, "\n")
	return content, true, nil
}

func ReadPEMFile(filepath string) (base64Str string, isSet bool, err error) {
	pemData, isSet, err := ReadFromFile(filepath)
	if err != nil {
		return "", false, fmt.Errorf("reading file: %w", err)
	} else if !isSet {
		return "", false, nil
	}

	base64Str, err = extract.PEM([]byte(pemData))
	if err != nil {
		return "", false, fmt.Errorf("extracting base64 encoded data from PEM content: %w", err)
	}

	return base64Str, true, nil
}
