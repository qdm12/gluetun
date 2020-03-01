package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetUser obtains the user to use to connect to the VPN servers
func (p *paramsReader) GetUser() (s string, err error) {
	defer func() {
		unsetenvErr := p.unsetEnv("USER")
		if err == nil {
			err = unsetenvErr
		}
	}()
	return p.envParams.GetEnv("USER", libparams.CaseSensitiveValue(), libparams.Compulsory())
}

// GetPassword obtains the password to use to connect to the VPN servers
func (p *paramsReader) GetPassword() (s string, err error) {
	defer func() {
		unsetenvErr := p.unsetEnv("PASSWORD")
		if err == nil {
			err = unsetenvErr
		}
	}()
	return p.envParams.GetEnv("PASSWORD", libparams.CaseSensitiveValue(), libparams.Compulsory())
}

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL
func (p *paramsReader) GetNetworkProtocol() (protocol models.NetworkProtocol, err error) {
	s, err := p.envParams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return models.NetworkProtocol(s), err
}

// GetOpenVPNVerbosity obtains the verbosity level for verbosity between 0 and 6
// from the environment variable OPENVPN_VERBOSITY
func (p *paramsReader) GetOpenVPNVerbosity() (verbosity int, err error) {
	return p.envParams.GetEnvIntRange("OPENVPN_VERBOSITY", 0, 6, libparams.Default("1"))
}
