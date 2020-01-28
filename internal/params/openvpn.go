package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetNonRoot obtains if openvpn should be run without root from the
// environment variable NONROOT
func (p *paramsReader) GetNonRoot() (nonRoot bool, err error) {
	return p.envParams.GetYesNo("NONROOT", libparams.Default("yes"))
}

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL
func (p *paramsReader) GetNetworkProtocol() (protocol constants.NetworkProtocol, err error) {
	s, err := p.envParams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return constants.NetworkProtocol(s), err
}
