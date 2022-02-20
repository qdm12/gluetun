package parse

import (
	"fmt"
)

func ExtractCert(b []byte) (certData string, err error) {
	certData, err = extractPEM(b, "CERTIFICATE")
	if err != nil {
		return "", fmt.Errorf("cannot extract PEM data: %w", err)
	}

	return certData, nil
}
