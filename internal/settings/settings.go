package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// Settings contains all settings for the program to run
type Settings struct {
	VPNSP       models.VPNProvider
	OpenVPN     OpenVPN
	PIA         PIA
	Mullvad     Mullvad
	Windscribe  Windscribe
	Surfshark   Surfshark
	System      System
	DNS         DNS
	Firewall    Firewall
	TinyProxy   TinyProxy
	ShadowSocks ShadowSocks
}

func (s *Settings) String() string {
	var vpnServiceProviderSettings string
	switch s.VPNSP {
	case constants.PrivateInternetAccess:
		vpnServiceProviderSettings = s.PIA.String()
	case constants.Mullvad:
		vpnServiceProviderSettings = s.Mullvad.String()
	case constants.Windscribe:
		vpnServiceProviderSettings = s.Windscribe.String()
	case constants.Surfshark:
		vpnServiceProviderSettings = s.Surfshark.String()
	}
	return strings.Join([]string{
		"Settings summary below:",
		s.OpenVPN.String(),
		vpnServiceProviderSettings,
		s.System.String(),
		s.DNS.String(),
		s.Firewall.String(),
		s.TinyProxy.String(),
		s.ShadowSocks.String(),
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
	settings.OpenVPN, err = GetOpenVPNSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	switch settings.VPNSP {
	case constants.PrivateInternetAccess:
		settings.PIA, err = GetPIASettings(paramsReader)
	case constants.Mullvad:
		settings.Mullvad, err = GetMullvadSettings(paramsReader)
	case constants.Windscribe:
		settings.Windscribe, err = GetWindscribeSettings(paramsReader, settings.OpenVPN.NetworkProtocol)
	case constants.Surfshark:
		settings.Surfshark, err = GetSurfsharkSettings(paramsReader)
	default:
		err = fmt.Errorf("VPN service provider %q is not valid", settings.VPNSP)
	}
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
	return settings, nil
}
