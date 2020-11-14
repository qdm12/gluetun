package params

import (
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetCyberghostGroup obtains the server group for the Cyberghost server from the
// environment variable CYBERGHOST_GROUP.
func (p *reader) GetCyberghostGroup() (group string, err error) {
	s, err := p.envParams.GetValueIfInside("CYBERGHOST_GROUP",
		constants.CyberghostGroupChoices(), libparams.Default("Premium UDP Europe"))
	return s, err
}

// GetCyberghostRegions obtains the country names for the Cyberghost servers from the
// environment variable REGION.
func (p *reader) GetCyberghostRegions() (regions []string, err error) {
	return p.envParams.GetCSVInPossibilities("REGION", constants.CyberghostRegionChoices())
}

// GetCyberghostClientKey obtains the one line client key to use for openvpn from the
// environment variable CLIENT_KEY or from the file at /gluetun/client.key.
func (p *reader) GetCyberghostClientKey() (clientKey string, err error) {
	clientKey, err = p.envParams.GetEnv("CLIENT_KEY", libparams.CaseSensitiveValue())
	if err != nil {
		return "", err
	} else if len(clientKey) > 0 {
		return clientKey, nil
	}
	content, err := p.fileManager.ReadFile(string(constants.ClientKey))
	if err != nil {
		return "", err
	}
	return extractClientKey(content)
}

func extractClientKey(b []byte) (key string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", fmt.Errorf("cannot decode PEM block from client key")
	}
	parsedBytes := pem.EncodeToMemory(pemBlock)
	s := string(parsedBytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN PRIVATE KEY-----")
	s = strings.TrimSuffix(s, "-----END PRIVATE KEY-----")
	return s, nil
}

// GetCyberghostClientCertificate obtains the client certificate to use for openvpn from the
// file at /gluetun/client.crt.
func (p *reader) GetCyberghostClientCertificate() (clientCertificate string, err error) {
	content, err := p.fileManager.ReadFile(string(constants.ClientCertificate))
	if err != nil {
		return "", err
	}
	return extractClientCertificate(content)
}

func extractClientCertificate(b []byte) (certificate string, err error) {
	pemBlock, _ := pem.Decode(b)
	if pemBlock == nil {
		return "", fmt.Errorf("cannot decode PEM block from client certificate")
	}
	parsedBytes := pem.EncodeToMemory(pemBlock)
	s := string(parsedBytes)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.TrimPrefix(s, "-----BEGIN CERTIFICATE-----")
	s = strings.TrimSuffix(s, "-----END CERTIFICATE-----")
	return s, nil
}
