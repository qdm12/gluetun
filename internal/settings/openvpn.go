package settings

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// OpenVPN contains settings to configure the OpenVPN client
type OpenVPN struct {
	NetworkProtocol models.NetworkProtocol
	Verbosity       int
	Root            bool
	TargetIP        net.IP
	Cipher          string
	Auth            string
}

// GetOpenVPNSettings obtains the OpenVPN settings using the params functions
func GetOpenVPNSettings(paramsReader params.Reader) (settings OpenVPN, err error) {
	settings.NetworkProtocol, err = paramsReader.GetNetworkProtocol()
	if err != nil {
		return settings, err
	}
	settings.Verbosity, err = paramsReader.GetOpenVPNVerbosity()
	if err != nil {
		return settings, err
	}
	settings.Root, err = paramsReader.GetOpenVPNRoot()
	if err != nil {
		return settings, err
	}
	settings.TargetIP, err = paramsReader.GetTargetIP()
	if err != nil {
		return settings, err
	}
	settings.Cipher, err = paramsReader.GetOpenVPNCipher()
	if err != nil {
		return settings, err
	}
	settings.Auth, err = paramsReader.GetOpenVPNAuth()
	if err != nil {
		return settings, err
	}
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
		"Target IP address: " + o.TargetIP.String(),
		"Custom cipher: " + o.Cipher,
		"Custom auth algorithm: " + o.Auth,
	}
	return strings.Join(settingsList, "\n|--")
}
