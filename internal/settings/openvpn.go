package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// OpenVPN contains settings to configure the OpenVPN client
type OpenVPN struct {
	NetworkProtocol models.NetworkProtocol
	Verbosity       int
}

// GetOpenVPNSettings obtains the OpenVPN settings using the params functions
func GetOpenVPNSettings(params params.ParamsReader) (settings OpenVPN, err error) {
	settings.NetworkProtocol, err = params.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.Verbosity, err = params.GetOpenVPNVerbosity()
	if err != nil {
		return settings, err
	}
	return settings, nil
}

func (o *OpenVPN) String() string {
	settingsList := []string{
		"OpenVPN settings:",
		"Network protocol: " + string(o.NetworkProtocol),
		"Verbosity level: " + fmt.Sprintf("%d", o.Verbosity),
	}
	return strings.Join(settingsList, "\n|--")
}
