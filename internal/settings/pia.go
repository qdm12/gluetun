package settings

import (
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
		return "Port forwarding: on, saved in " + string(p.Filepath)
	}
	return "Port forwarding: off"
}

func (p *PIA) String() string {
	settingsList := []string{
		"Region: " + string(p.Region),
		"Encryption: " + string(p.Encryption),
		"Port forwarding: " + p.PortForwarding.String(),
	}
	return "PIA settings:\n" + strings.Join(settingsList, "\n |--")
}

// GetPIASettings obtains PIA settings from environment variables using the params package.
func GetPIASettings(params params.ParamsReader) (settings PIA, err error) {
	settings.User, err = params.GetUser()
	if err != nil {
		return settings, err
	}
	settings.Password, err = params.GetPassword()
	if err != nil {
		return settings, err
	}
	settings.Encryption, err = params.GetPIAEncryption()
	if err != nil {
		return settings, err
	}
	settings.Region, err = params.GetPIARegion()
	if err != nil {
		return settings, err
	}
	settings.PortForwarding.Enabled, err = params.GetPortForwarding()
	if err != nil {
		return settings, err
	}
	if settings.PortForwarding.Enabled {
		settings.PortForwarding.Filepath, err = params.GetPortForwardingStatusFilepath()
		if err != nil {
			return settings, err
		}
	}
	return settings, nil
}
