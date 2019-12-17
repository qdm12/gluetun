package settings

import "github.com/qdm12/private-internet-access-docker/internal/constants"
import "github.com/qdm12/private-internet-access-docker/internal/params"

import "strings"

// OpenVPN contains settings to configure the OpenVPN client
type OpenVPN struct {
	NonRoot         bool
	NetworkProtocol constants.NetworkProtocol
}

// GetOpenVPNSettings obtains the OpenVPN settings using the params functions
func GetOpenVPNSettings() (settings OpenVPN, err error) {
	nonRoot, err := params.GetNonRoot()
	if err != nil {
		return settings, err
	}
	settings.NonRoot = nonRoot
	networkProtocol, err := params.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.NetworkProtocol = networkProtocol
	return settings, nil
}

func (o *OpenVPN) String() string {
	nonRootStr := "on"
	if !o.NonRoot {
		nonRootStr = "off"
	}
	settingsList := []string{
		"Running without root privileges: " + nonRootStr,
		"Network protocol: " + string(o.NetworkProtocol),
	}
	return "OpenVPN settings:\n" + strings.Join(settingsList, "\n|--")
}
