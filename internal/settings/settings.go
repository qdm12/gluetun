package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Settings contains all settings for the program to run
type Settings struct {
	VPNSP       string
	OpenVPN     OpenVPN
	PIA         PIA
	Mullvad     Mullvad
	Windscribe  Windscribe
	DNS         DNS
	Firewall    Firewall
	TinyProxy   TinyProxy
	ShadowSocks ShadowSocks
}

func (s *Settings) String() string {
	var vpnServiceProvider string
	switch s.VPNSP {
	case "pia":
		vpnServiceProvider = s.PIA.String()
	case "mullvad":
		vpnServiceProvider = s.Mullvad.String()
	case "windscribe":
		vpnServiceProvider = s.Windscribe.String()
	}
	return strings.Join([]string{
		"Settings summary below:",
		s.OpenVPN.String(),
		vpnServiceProvider,
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
	settings.VPNSP, err = params.GetVPNSP()
	if err != nil {
		return settings, err
	}
	settings.OpenVPN, err = GetOpenVPNSettings(params)
	if err != nil {
		return settings, err
	}
	switch settings.VPNSP {
	case "pia":
		settings.PIA, err = GetPIASettings(params)
	case "mullvad":
		settings.Mullvad, err = GetMullvadSettings(params)
	case "windscribe":
		settings.Windscribe, err = GetWindscribeSettings(params)
	default:
		err = fmt.Errorf("VPN service provider %q is not valid", settings.VPNSP)
	}
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
