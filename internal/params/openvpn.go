package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL
func (p *paramsReader) GetNetworkProtocol() (protocol models.NetworkProtocol, err error) {
	s, err := p.envParams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return models.NetworkProtocol(s), err
}
