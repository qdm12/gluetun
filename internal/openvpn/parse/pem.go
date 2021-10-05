package parse

import (
	"encoding/pem"
	"errors"
	"strings"
)

var (
	errPEMDecode = errors.New("cannot decode PEM encoded block")
)

func extractPEM(b []byte, name string) (encodedData string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", errPEMDecode
	}

	encodedBytes := pem.EncodeToMemory(pemBlock)
	encodedData = string(encodedBytes)
	encodedData = strings.ReplaceAll(encodedData, "\n", "")
	encodedData = strings.TrimPrefix(encodedData, "-----BEGIN "+name+"-----")
	encodedData = strings.TrimSuffix(encodedData, "-----END "+name+"-----")
	return encodedData, nil
}
