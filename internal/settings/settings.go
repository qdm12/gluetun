package settings

import (
	"fmt"
	"strings"
	"time"

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
	PublicIPPeriod     time.Duration
	UpdaterPeriod      time.Duration
	VersionInformation bool
	ControlServer      ControlServer
}

func (s *Settings) String() string {
	versionInformation := disabled
	if s.VersionInformation {
		versionInformation = enabled
	}
	updaterLine := "Updater: disabled"
	if s.UpdaterPeriod > 0 {
		updaterLine = fmt.Sprintf("Updater period: %s", s.UpdaterPeriod)
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
		"Public IP check period: " + s.PublicIPPeriod.String(), // TODO print disabled if 0
		"Version information: " + versionInformation,
		updaterLine,
		"", // new line at the end
	}, "\n")
}

// GetAllSettings obtains all settings for the program and returns an error as soon
// as an error is encountered reading them.
func GetAllSettings(paramsReader params.Reader) (settings Settings, err error) {
	settings.VPNSP, err = paramsReader.GetVPNSP()
	if err != nil {
		return settings, err
	}
	settings.OpenVPN, err = GetOpenVPNSettings(paramsReader, settings.VPNSP)
	if err != nil {
		return settings, err
	}
	settings.DNS, err = GetDNSSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	settings.Firewall, err = GetFirewallSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	settings.HTTPProxy, err = GetHTTPProxySettings(paramsReader)
	if err != nil {
		return settings, err
	}
	settings.ShadowSocks, err = GetShadowSocksSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	settings.System, err = GetSystemSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	settings.PublicIPPeriod, err = paramsReader.GetPublicIPPeriod()
	if err != nil {
		return settings, err
	}
	settings.VersionInformation, err = paramsReader.GetVersionInformation()
	if err != nil {
		return settings, err
	}
	settings.UpdaterPeriod, err = paramsReader.GetUpdaterPeriod()
	if err != nil {
		return settings, err
	}
	settings.ControlServer, err = GetControlServerSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
