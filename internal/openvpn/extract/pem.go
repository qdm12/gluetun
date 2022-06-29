package extract

import (
	"encoding/pem"
	"errors"
	"regexp"
	"strings"
)

var (
	errPEMDecode = errors.New("cannot decode PEM encoded block")
)

var (
	regexPEMBegin = regexp.MustCompile(`-----BEGIN [A-Za-z ]+-----`)
	regexPEMEnd   = regexp.MustCompile(`-----END [A-Za-z ]+-----`)
)

func PEM(b []byte) (encodedData string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", errPEMDecode
	}

	encodedBytes := pem.EncodeToMemory(pemBlock)
	encodedData = string(encodedBytes)
	encodedData = strings.ReplaceAll(encodedData, "\n", "")
	beginPrefix := regexPEMBegin.FindString(encodedData)
	encodedData = strings.TrimPrefix(encodedData, beginPrefix)
	endPrefix := regexPEMEnd.FindString(encodedData)
	encodedData = strings.TrimSuffix(encodedData, endPrefix)
	return encodedData, nil
}
