package params

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/models"
	libparams "github.com/qdm12/golibs/params"
)

// GetUser obtains the user to use to connect to the VPN servers
func (r *reader) GetUser() (s string, err error) {
	defer func() {
		unsetenvErr := r.unsetEnv("USER")
		if err == nil {
			err = unsetenvErr
		}
	}()
	return r.envParams.GetEnv("USER", libparams.CaseSensitiveValue(), libparams.Compulsory())
}

// GetPassword obtains the password to use to connect to the VPN servers
func (r *reader) GetPassword(required bool) (s string, err error) {
	defer func() {
		unsetenvErr := r.unsetEnv("PASSWORD")
		if err == nil {
			err = unsetenvErr
		}
	}()
	options := []libparams.GetEnvSetter{libparams.CaseSensitiveValue()}
	if required {
		options = append(options, libparams.Compulsory())
	}
	return r.envParams.GetEnv("PASSWORD", options...)
}

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL
func (r *reader) GetNetworkProtocol() (protocol models.NetworkProtocol, err error) {
	s, err := r.envParams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, libparams.Default("udp"))
	return models.NetworkProtocol(s), err
}

// GetOpenVPNVerbosity obtains the verbosity level for verbosity between 0 and 6
// from the environment variable OPENVPN_VERBOSITY
func (r *reader) GetOpenVPNVerbosity() (verbosity int, err error) {
	return r.envParams.GetEnvIntRange("OPENVPN_VERBOSITY", 0, 6, libparams.Default("1"))
}

// GetOpenVPNRoot obtains if openvpn should be run as root
// from the environment variable OPENVPN_ROOT
func (r *reader) GetOpenVPNRoot() (root bool, err error) {
	return r.envParams.GetYesNo("OPENVPN_ROOT", libparams.Default("no"))
}

// GetTargetIP obtains the IP address to choose from the list of IP addresses
// available for a particular region, from the environment variable
// OPENVPN_TARGET_IP
func (r *reader) GetTargetIP() (ip net.IP, err error) {
	s, err := r.envParams.GetEnv("OPENVPN_TARGET_IP")
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
func (r *reader) GetOpenVPNCipher() (cipher string, err error) {
	return r.envParams.GetEnv("OPENVPN_CIPHER")
}

// GetOpenVPNAuth obtains a custom auth algorithm to use with OpenVPN
// from the environment variable OPENVPN_AUTH
func (r *reader) GetOpenVPNAuth() (auth string, err error) {
	return r.envParams.GetEnv("OPENVPN_AUTH")
}
