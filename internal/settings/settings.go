package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/params"
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
		"Settings summary below:",
		s.OpenVPN.String(),
		s.PIA.String(),
		s.DNS.String(),
		s.Firewall.String(),
		s.TinyProxy.String(),
		s.ShadowSocks.String(),
		"", // new line at the end
	}, "\n")
}

// GetAllSettings obtains all settings for the program and returns an error as soon
// as an error is encountered reading them.
func GetAllSettings(params params.ParamsReader) (settings Settings, err error) {
	settings.OpenVPN, err = GetOpenVPNSettings(params)
	if err != nil {
		return settings, err
	}
	settings.PIA, err = GetPIASettings(params)
	if err != nil {
		return settings, err
	}
	settings.DNS, err = GetDNSSettings(params)
	if err != nil {
		return settings, err
	}
	settings.Firewall, err = GetFirewallSettings(params)
	if err != nil {
		return settings, err
	}
	settings.TinyProxy, err = GetTinyProxySettings(params)
	if err != nil {
		return settings, err
	}
	settings.ShadowSocks, err = GetShadowSocksSettings(params)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
