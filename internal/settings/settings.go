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
		switch settings.OpenVPN.Cipher {
		case "", "aes-128-cbc", "aes-256-cbc", "aes-128-gcm", "aes-256-gcm":
		default:
			return settings, fmt.Errorf("cipher %q is not supported by Private Internet Access", settings.OpenVPN.Cipher)
		}
		switch settings.OpenVPN.Auth {
		case "", "sha1", "sha256":
		default:
			return settings, fmt.Errorf("auth algorithm %q is not supported by Private Internet Access", settings.OpenVPN.Auth)
		}
		settings.PIA, err = GetPIASettings(paramsReader)
	case constants.Mullvad:
		switch settings.OpenVPN.Cipher {
		case "":
		default:
			return settings, fmt.Errorf("cipher %q is not supported by Mullvad", settings.OpenVPN.Cipher)
		}
		switch settings.OpenVPN.Auth {
		case "":
		default:
			return settings, fmt.Errorf("auth algorithm %q is not supported by Mullvad (not using auth at all)", settings.OpenVPN.Auth)
		}
		settings.Mullvad, err = GetMullvadSettings(paramsReader)
	case constants.Windscribe:
		switch settings.OpenVPN.Cipher {
		case "", "aes-256-cbc", "aes-256-gcm": // TODO check inside params getters
		default:
			return settings, fmt.Errorf("cipher %q is not supported by Windscribe", settings.OpenVPN.Cipher)
		}
		switch settings.OpenVPN.Auth {
		case "", "sha512":
		default:
			return settings, fmt.Errorf("auth algorithm %q is not supported by Windscribe", settings.OpenVPN.Auth)
		}
		settings.Windscribe, err = GetWindscribeSettings(paramsReader, settings.OpenVPN.NetworkProtocol)
	default:
		err = fmt.Errorf("VPN service provider %q is not valid", settings.VPNSP)
	}
	if err != nil {
		return settings, err
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
