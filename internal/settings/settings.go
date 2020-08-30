package settings

import (
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

const (
	enabled  = "enabled"
	disabled = "disabled"
)

// Settings contains all settings for the program to run
type Settings struct {
	VPNSP              models.VPNProvider
	OpenVPN            OpenVPN
	System             System
	DNS                DNS
	Firewall           Firewall
	TinyProxy          TinyProxy
	ShadowSocks        ShadowSocks
	PublicIPPeriod     time.Duration
	VersionInformation bool
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
		s.TinyProxy.String(),
		s.ShadowSocks.String(),
		"Public IP check period: " + s.PublicIPPeriod.String(),
		"Version information: " + versionInformation,
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
	settings.TinyProxy, err = GetTinyProxySettings(paramsReader)
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
	return settings, nil
}
