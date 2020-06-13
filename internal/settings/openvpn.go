package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// OpenVPN contains settings to configure the OpenVPN client
type OpenVPN struct {
	User      string
	Password  string
	Verbosity int
	Root      bool
	Cipher    string
	Auth      string
}

// GetOpenVPNSettings obtains the OpenVPN settings using the params functions
func GetOpenVPNSettings(paramsReader params.Reader) (settings OpenVPN, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	// Remove spaces in user ID to simplify user's life, thanks @JeordyR
	settings.User = strings.ReplaceAll(settings.User, " ", "")
	settings.Password, err = paramsReader.GetPassword()
	if err != nil {
		return settings, err
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
	return settings, nil
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
	}
	if len(o.Cipher) > 0 {
		settingsList = append(settingsList, "Custom cipher: "+o.Cipher)
	}
	if len(o.Auth) > 0 {
		settingsList = append(settingsList, "Custom auth algorithm: "+o.Auth)
	}
	return strings.Join(settingsList, "\n|--")
}
