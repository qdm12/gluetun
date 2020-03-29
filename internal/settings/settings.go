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
	UID         int
	GID         int
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
		fmt.Sprintf("|-- Using UID %d and GID %d", s.UID, s.GID),
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
		settings.PIA, err = GetPIASettings(params)
	case "mullvad":
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
		settings.Mullvad, err = GetMullvadSettings(params)
	case "windscribe":
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
		settings.Windscribe, err = GetWindscribeSettings(params, settings.OpenVPN.NetworkProtocol)
	default:
		err = fmt.Errorf("VPN service provider %q is not valid", settings.VPNSP)
	}
	if err != nil {
		return settings, err
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
	settings.UID, err = params.GetUID()
	if err != nil {
		return settings, err
	}
	settings.GID, err = params.GetGID()
	if err != nil {
		return settings, err
	}
	return settings, nil
}
