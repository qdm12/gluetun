package params

import (
	"fmt"

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
	s, err = p.envParams.GetEnv("USER")
	if err != nil {
		return "", err
	} else if len(s) == 0 {
		return s, fmt.Errorf("USER environment variable cannot be empty")
	}
	s = strings.Replace(s, ' ', '', -1)
	return s, nil
}

// GetPassword obtains the password to use to connect to the VPN servers
func (p *paramsReader) GetPassword() (s string, err error) {
	defer func() {
		unsetenvErr := p.unsetEnv("PASSWORD")
		if err == nil {
			err = unsetenvErr
		}
	}()
	s, err = p.envParams.GetEnv("PASSWORD")
	if err != nil {
		return "", err
	} else if len(s) == 0 {
		return s, fmt.Errorf("PASSWORD environment variable cannot be empty")
	}
	return s, nil
}

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL
func (p *paramsReader) GetNetworkProtocol() (protocol models.NetworkProtocol, err error) {
	s, err := p.envParams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return models.NetworkProtocol(s), err
}
