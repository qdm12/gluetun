package params

import (
	"fmt"
	"net"
	"strings"

	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

// GetUser obtains the user to use to connect to the VPN servers
func (p *reader) GetUser() (s string, err error) {
	defer func() {
		unsetenvErr := p.unsetEnv("USER")
		if err == nil {
			err = unsetenvErr
		}
	}()
	return p.envParams.GetEnv("USER", libparams.CaseSensitiveValue(), libparams.Compulsory())
}

// GetPassword obtains the password to use to connect to the VPN servers
func (p *reader) GetPassword() (s string, err error) {
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
func (p *reader) GetNetworkProtocol() (protocol models.NetworkProtocol, err error) {
	s, err := p.envParams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return models.NetworkProtocol(s), err
}

// GetOpenVPNVerbosity obtains the verbosity level for verbosity between 0 and 6
// from the environment variable OPENVPN_VERBOSITY
func (p *reader) GetOpenVPNVerbosity() (verbosity int, err error) {
	return p.envParams.GetEnvIntRange("OPENVPN_VERBOSITY", 0, 6, libparams.Default("1"))
}

// GetOpenVPNRoot obtains if openvpn should be run as root
// from the environment variable OPENVPN_ROOT
func (p *reader) GetOpenVPNRoot() (root bool, err error) {
	return p.envParams.GetYesNo("OPENVPN_ROOT", libparams.Default("no"))
}

// GetTargetIP obtains the IP address to choose from the list of IP addresses
// available for a particular region, from the environment variable
// OPENVPN_TARGET_IP
func (p *reader) GetTargetIP() (ip net.IP, err error) {
	s, err := p.envParams.GetEnv("OPENVPN_TARGET_IP")
	if len(s) == 0 {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	ip = net.ParseIP(s)
	if ip == nil {
		return nil, fmt.Errorf("target IP address %q is not valid", s)
	}
	return ip, nil
}

// GetOpenVPNCipher obtains a custom cipher to use with OpenVPN
// from the environment variable OPENVPN_CIPHER
func (p *reader) GetOpenVPNCipher() (cipher string, err error) {
	cipher, err = p.envParams.GetEnv("OPENVPN_CIPHER")
	return strings.ToLower(cipher), err
}

// GetOpenVPNAuth obtains a custom auth algorithm to use with OpenVPN
// from the environment variable OPENVPN_AUTH
func (p *reader) GetOpenVPNAuth() (auth string, err error) {
	auth, err = p.envParams.GetEnv("OPENVPN_AUTH")
	return strings.ToLower(auth), err
}
