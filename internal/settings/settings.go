package settings

import (
	"strings"

	libparams "github.com/qdm12/golibs/params"
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
func GetAllSettings(envParams libparams.EnvParams) (settings Settings, err error) {
	settings.OpenVPN, err = GetOpenVPNSettings(envParams)
	if err != nil {
		return settings, err
	}
	settings.PIA, err = GetPIASettings(envParams)
	if err != nil {
		return settings, err
	}
	settings.DNS, err = GetDNSSettings(envParams)
	if err != nil {
		return settings, err
	}
	settings.Firewall, err = GetFirewallSettings(envParams)
	if err != nil {
		return settings, err
	}
	settings.TinyProxy, err = GetTinyProxySettings(envParams)
	if err != nil {
		return settings, err
	}
	settings.ShadowSocks, err = GetShadowSocksSettings(envParams)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
