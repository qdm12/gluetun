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
	Provider    models.ProviderSettings
	System      System
	DNS         DNS
	Firewall    Firewall
	TinyProxy   TinyProxy
	ShadowSocks ShadowSocks
}

func (s *Settings) String() string {
	return strings.Join([]string{
		"Settings summary below:",
		s.OpenVPN.String(),
		s.Provider.String(),
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
	switch settings.VPNSP {
	case constants.PrivateInternetAccess:
		settings.Provider, err = GetPIASettings(paramsReader)
	case constants.Mullvad:
		settings.Provider, err = GetMullvadSettings(paramsReader)
	case constants.Windscribe:
		settings.Provider, err = GetWindscribeSettings(paramsReader)
	case constants.Surfshark:
		settings.Provider, err = GetSurfsharkSettings(paramsReader)
	case constants.Cyberghost:
		settings.Provider, err = GetCyberghostSettings(paramsReader)
	default:
		err = fmt.Errorf("VPN service provider %q is not valid", settings.VPNSP)
	}
	if err != nil {
		return settings, err
	}
	settings.OpenVPN, err = GetOpenVPNSettings(paramsReader)
	if err != nil {
		return settings, err
	}
	if settings.VPNSP == constants.Mullvad {
		settings.OpenVPN.Password = "m"
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
