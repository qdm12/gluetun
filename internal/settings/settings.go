package settings

import (
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

const (
	enabled  = "enabled"
	disabled = "disabled"
)

// Settings contains all settings for the program to run.
type Settings struct {
	VPNSP              models.VPNProvider
	OpenVPN            OpenVPN
	System             System
	DNS                DNS
	Firewall           Firewall
	HTTPProxy          HTTPProxy
	ShadowSocks        ShadowSocks
	Updater            Updater
	PublicIP           PublicIP
	VersionInformation bool
	ControlServer      ControlServer
}

func (s *Settings) String() string {
	versionInformation := disabled
	if s.VersionInformation {
		versionInformation = enabled
	}
	return strings.Join([]string{
		"Settings summary below:",
		s.OpenVPN.String(),
		s.System.String(),
		s.DNS.String(),
		s.Firewall.String(),
		s.HTTPProxy.String(),
		s.ShadowSocks.String(),
		s.ControlServer.String(),
		s.Updater.String(),
		s.PublicIP.String(),
		"Version information: " + versionInformation,
		"", // new line at the end
	}, "\n")
}

// GetAllSettings obtains all settings for the program and returns an error as soon
// as an error is encountered reading them.
func GetAllSettings(paramsReader params.Reader) (settings Settings, warnings []string, err error) {
	settings.VPNSP, err = paramsReader.GetVPNSP()
	if err != nil {
		return settings, nil, err
	}
	settings.OpenVPN, err = GetOpenVPNSettings(paramsReader, settings.VPNSP)
	if err != nil {
		return settings, nil, err
	}
	settings.DNS, err = GetDNSSettings(paramsReader)
	if err != nil {
		return settings, nil, err
	}
	settings.Firewall, err = GetFirewallSettings(paramsReader)
	if err != nil {
		return settings, nil, err
	}
	settings.System, err = GetSystemSettings(paramsReader)
	if err != nil {
		return settings, nil, err
	}
	settings.PublicIP, err = getPublicIPSettings(paramsReader)
	if err != nil {
		return settings, nil, err
	}
	settings.VersionInformation, err = paramsReader.GetVersionInformation()
	if err != nil {
		return settings, nil, err
	}
	settings.Updater, err = GetUpdaterSettings(paramsReader)
	if err != nil {
		return settings, nil, err
	}

	var warning string
	settings.HTTPProxy, warning, err = GetHTTPProxySettings(paramsReader)
	if warning != "" {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return settings, warnings, err
	}

	settings.ShadowSocks, warning, err = GetShadowSocksSettings(paramsReader)
	if warning != "" {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return settings, warnings, err
	}

	settings.ControlServer, warning, err = GetControlServerSettings(paramsReader)
	if warning != "" {
		warnings = append(warnings, warning)
	}
	if err != nil {
		return settings, warnings, err
	}

	return settings, warnings, nil
}
