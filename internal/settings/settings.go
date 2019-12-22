package settings

import (
	"strings"
)

// Settings contains all settings for the program to run
type Settings struct {
	OpenVPN     OpenVPN
	PIA         PIA
	DNS         DNS
	Firewall    Firewall
	TinyProxy   TinyProxy
	ShadowSocks ShadowSocks
}

func (s *Settings) String() string {
	return strings.Join([]string{
		s.OpenVPN.String(),
		s.PIA.String(),
		s.DNS.String(),
		s.Firewall.String(),
		s.TinyProxy.String(),
		s.ShadowSocks.String(),
	}, "\n")
}

// GetAllSettings obtains all settings for the program and returns an error as soon
// as an error is encountered reading them.
func GetAllSettings() (settings Settings, err error) {
	settings.OpenVPN, err = GetOpenVPNSettings()
	if err != nil {
		return settings, err
	}
	settings.PIA, err = GetPIASettings()
	if err != nil {
		return settings, err
	}
	settings.DNS, err = GetDNSSettings()
	if err != nil {
		return settings, err
	}
	settings.Firewall, err = GetFirewallSettings()
	if err != nil {
		return settings, err
	}
	settings.TinyProxy, err = GetTinyProxySettings()
	if err != nil {
		return settings, err
	}
	settings.ShadowSocks, err = GetShadowSocksSettings()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
