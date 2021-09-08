package parse

import (
	"fmt"
)

func ExtractCert(b []byte) (certData string, err error) {
	certData, err = extractPEM(b, "CERTIFICATE")
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrExtractPEM, err)
	}

	return certData, nil
}
