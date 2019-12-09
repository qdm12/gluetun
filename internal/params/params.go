package params

import (
	"fmt"

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
	s := libparams.GetEnv("PROTOCOL", "tcp")
	if s == "tcp" {
		return constants.TCP, nil
	} else if s == "udp" {
		return constants.UDP, nil
	}
	return 0, fmt.Errorf("PROTOCOL can only be \"tcp\" or \"udp\"")
}
