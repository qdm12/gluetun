package params

import (
	libparams "github.com/qdm12/golibs/params"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

// GetNonRoot obtains if openvpn should be run without root from the
// environment variable NONROOT
func GetNonRoot() (nonRoot bool, err error) {
	return libparams.GetYesNo("NONROOT", true)
}

// GetNetworkProtocol obtains the network protocol to use to connect to the
// VPN servers from the environment variable PROTOCOL
func GetNetworkProtocol() (protocol constants.NetworkProtocol, err error) {
	s, err := libparams.GetValueIfInside("PROTOCOL", []string{"tcp", "udp"}, false, "udp")
	return constants.NetworkProtocol(s), err
}
