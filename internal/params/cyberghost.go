package params

import (
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	libparams "github.com/qdm12/golibs/params"
)

// GetCyberghostGroup obtains the server group for the Cyberghost server from the
// environment variable CYBERGHOST_GROUP
func (p *reader) GetCyberghostGroup() (group string, err error) {
	s, err := p.envParams.GetValueIfInside("CYBERGHOST_GROUP", constants.CyberghostGroupChoices(), libparams.Default("Premium UDP Europe"))
	return s, err
}

// GetCyberghostRegion obtains the country name for the Cyberghost server from the
// environment variable REGION
func (p *reader) GetCyberghostRegion() (region string, err error) {
	choices := append(constants.CyberghostRegionChoices(), "")
	s, err := p.envParams.GetValueIfInside("REGION", choices)
	return s, err
}

// GetCyberghostClientKey obtains the one line client key to use for openvpn from the
// environment variable CLIENT_KEY
func (p *reader) GetCyberghostClientKey() (clientKey string, err error) {
	clientKey, err = p.envParams.GetEnv("CLIENT_KEY", libparams.CaseSensitiveValue())
	if err != nil {
		return "", err
	} else if len(clientKey) > 0 {
		return clientKey, nil
	}
	content, err := p.fileManager.ReadFile("/files/client.key")
	if err != nil {
		return "", err
	}
	s := string(content)
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s, nil
}
