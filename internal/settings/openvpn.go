package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/params"
)

// OpenVPN contains settings to configure the OpenVPN client
type OpenVPN struct {
	User      string                  `json:"user"`
	Password  string                  `json:"-"`
	Verbosity int                     `json:"verbosity"`
	Root      bool                    `json:"runAsRoot"`
	Cipher    string                  `json:"cipher"`
	Auth      string                  `json:"auth"`
	Provider  models.ProviderSettings `json:"provider"`
}

// GetOpenVPNSettings obtains the OpenVPN settings using the params functions
func GetOpenVPNSettings(paramsReader params.Reader, vpnProvider models.VPNProvider) (settings OpenVPN, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	settings.User = strings.ReplaceAll(settings.User, " ", "")
	isMullvad := vpnProvider == constants.Mullvad
	settings.Password, err = paramsReader.GetPassword(!isMullvad)
	if err != nil {
		return settings, err
	} else if isMullvad {
		settings.Password = "m"
	}
	settings.Verbosity, err = paramsReader.GetOpenVPNVerbosity()
	if err != nil {
		return settings, err
	}
	settings.Root, err = paramsReader.GetOpenVPNRoot()
	if err != nil {
		return settings, err
	}
	settings.Cipher, err = paramsReader.GetOpenVPNCipher()
	if err != nil {
		return settings, err
	}
	settings.Auth, err = paramsReader.GetOpenVPNAuth()
	if err != nil {
		return settings, err
	}
	switch vpnProvider {
	case constants.PrivateInternetAccess:
		settings.Provider, err = GetPIASettings(paramsReader)
	case constants.PrivateInternetAccessOld:
		settings.Provider, err = GetPIAOldSettings(paramsReader)
	case constants.Mullvad:
		settings.Provider, err = GetMullvadSettings(paramsReader)
	case constants.Windscribe:
		settings.Provider, err = GetWindscribeSettings(paramsReader)
	case constants.Surfshark:
		settings.Provider, err = GetSurfsharkSettings(paramsReader)
	case constants.Cyberghost:
		settings.Provider, err = GetCyberghostSettings(paramsReader)
	case constants.Vyprvpn:
		settings.Provider, err = GetVyprvpnSettings(paramsReader)
	case constants.Nordvpn:
		settings.Provider, err = GetNordvpnSettings(paramsReader)
	case constants.Purevpn:
		settings.Provider, err = GetPurevpnSettings(paramsReader)
	default:
		err = fmt.Errorf("VPN service provider %q is not valid", vpnProvider)
	}
	return settings, err
}

func (o *OpenVPN) String() string {
	runAsRoot := "no"
	if o.Root {
		runAsRoot = "yes"
	}
	settingsList := []string{
		"OpenVPN settings:",
		"User: [redacted]",
		"Password: [redacted]",
		"Verbosity level: " + fmt.Sprintf("%d", o.Verbosity),
		"Run as root: " + runAsRoot,
		o.Provider.String(),
	}
	if len(o.Cipher) > 0 {
		settingsList = append(settingsList, "Custom cipher: "+o.Cipher)
	}
	if len(o.Auth) > 0 {
		settingsList = append(settingsList, "Custom auth algorithm: "+o.Auth)
	}
	return strings.Join(settingsList, "\n|--")
}
