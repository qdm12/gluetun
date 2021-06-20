package configuration

import (
	"encoding/pem"
	"errors"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

func readClientKey(r reader) (clientKey string, err error) {
	b, err := r.getFromFileOrSecretFile("OPENVPN_CLIENTKEY", constants.ClientKey)
	if err != nil {
		return "", err
	}
	return extractClientKey(b)
}

var errDecodePEMBlockClientKey = errors.New("cannot decode PEM block from client key")

func extractClientKey(b []byte) (key string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", errDecodePEMBlockClientKey
	}
	parsedBytes := pem.EncodeToMemory(pemBlock)
	s := string(parsedBytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN PRIVATE KEY-----")
	s = strings.TrimSuffix(s, "-----END PRIVATE KEY-----")
	return s, nil
}

func readClientCertificate(r reader) (clientCertificate string, err error) {
	b, err := r.getFromFileOrSecretFile("OPENVPN_CLIENTCRT", constants.ClientCertificate)
	if err != nil {
		return "", err
	}
	return extractClientCertificate(b)
}

var errDecodePEMBlockClientCert = errors.New("cannot decode PEM block from client certificate")

func extractClientCertificate(b []byte) (certificate string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", errDecodePEMBlockClientCert
	}
	parsedBytes := pem.EncodeToMemory(pemBlock)
	s := string(parsedBytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN CERTIFICATE-----")
	s = strings.TrimSuffix(s, "-----END CERTIFICATE-----")
	return s, nil
}
