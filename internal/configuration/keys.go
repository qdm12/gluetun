package configuration

import (
	"errors"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/openvpn/parse"
)

var (
	errClientCert = errors.New("cannot read client certificate")
	errClientKey  = errors.New("cannot read client key")
)

func readClientKey(r reader) (clientKey string, err error) {
	b, err := r.getFromFileOrSecretFile("OPENVPN_CLIENTKEY", constants.ClientKey)
	if err != nil {
		return "", err
	}
	return parse.ExtractPrivateKey(b)
}

func readClientCertificate(r reader) (clientCertificate string, err error) {
	b, err := r.getFromFileOrSecretFile("OPENVPN_CLIENTCRT", constants.ClientCertificate)
	if err != nil {
		return "", err
	}
	return parse.ExtractCert(b)
}
