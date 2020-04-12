package settings

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// PIA contains the settings to connect to a PIA server
type PIA struct {
	User           string
	Password       string
	Encryption     models.PIAEncryption
	Region         models.PIARegion
	PortForwarding PortForwarding
}

// PortForwarding contains settings for port forwarding
type PortForwarding struct {
	Enabled  bool
	Filepath models.Filepath
}

func (p *PortForwarding) String() string {
	if p.Enabled {
		return fmt.Sprintf("on, saved in %s", p.Filepath)
	}
	return "off"
}

func (p *PIA) String() string {
	settingsList := []string{
		"PIA settings:",
		"User: [redacted]",
		"Password: [redacted]",
		"Region: " + string(p.Region),
		"Encryption: " + string(p.Encryption),
		"Port forwarding: " + p.PortForwarding.String(),
	}
	return strings.Join(settingsList, "\n |--")
}

// GetPIASettings obtains PIA settings from environment variables using the params package.
func GetPIASettings(paramsReader params.Reader) (settings PIA, err error) {
	settings.User, err = paramsReader.GetUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = paramsReader.GetPassword()
	if err != nil {
		return settings, err
	}
	settings.Encryption, err = paramsReader.GetPIAEncryption()
	if err != nil {
		return settings, err
	}
	settings.Region, err = paramsReader.GetPIARegion()
	if err != nil {
		return settings, err
	}
	settings.PortForwarding.Enabled, err = paramsReader.GetPortForwarding()
	if err != nil {
		return settings, err
	}
	if settings.PortForwarding.Enabled {
		settings.PortForwarding.Filepath, err = paramsReader.GetPortForwardingStatusFilepath()
		if err != nil {
			return settings, err
		}
	}
	return settings, nil
}
