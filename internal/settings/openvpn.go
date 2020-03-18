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
	Root            bool
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
	settings.Root, err = params.GetOpenVPNRoot()
	return settings, nil
}

func (o *OpenVPN) String() string {
	runAsRoot := "no"
	if o.Root {
		runAsRoot = "yes"
	}
	settingsList := []string{
		"OpenVPN settings:",
		"Network protocol: " + string(o.NetworkProtocol),
		"Verbosity level: " + fmt.Sprintf("%d", o.Verbosity),
		"Run as root: " + runAsRoot,
	}
	return strings.Join(settingsList, "\n|--")
}
