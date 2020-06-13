package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetCyberghostGroup obtains the server group for the Cyberghost server from the
// environment variable CYBERGHOST_GROUP
func (p *reader) GetCyberghostGroup() (group string, err error) {
	s, err := p.envParams.GetValueIfInside("CYBERGHOST_GROUP", constants.CyberghostGroupChoices())
	return s, err
}

// GetCyberghostRegion obtains the country name for the Cyberghost server from the
// environment variable REGION
func (p *reader) GetCyberghostRegion() (region string, err error) {
	s, err := p.envParams.GetValueIfInside("REGION", constants.CyberghostRegionChoices())
	return s, err
}

// GetCyberghostClientKey obtains the one line client key to use for openvpn from the
// environment variable CLIENT_KEY
func (p *reader) GetCyberghostClientKey() (clientKey string, err error) {
	return p.envParams.GetEnv("CLIENT_KEY", libparams.Compulsory(), libparams.CaseSensitiveValue())
}
