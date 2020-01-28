package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// OpenVPN contains settings to configure the OpenVPN client
type OpenVPN struct {
	NonRoot         bool
	NetworkProtocol constants.NetworkProtocol
}

// GetOpenVPNSettings obtains the OpenVPN settings using the params functions
func GetOpenVPNSettings(params params.ParamsReader) (settings OpenVPN, err error) {
	settings.NonRoot, err = params.GetNonRoot()
	if err != nil {
		return settings, err
	}
	settings.NetworkProtocol, err = params.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
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
