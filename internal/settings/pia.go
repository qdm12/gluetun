package settings

import (
	"strings"

	libparams "github.com/qdm12/golibs/params"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/params"
)

// PIA contains the settings to connect to a PIA server
type PIA struct {
	User           string
	Password       string
	Encryption     constants.PIAEncryption
	Region         constants.PIARegion
	PortForwarding PortForwarding
}

// PortForwarding contains settings for port forwarding
type PortForwarding struct {
	Enabled  bool
	Filepath string
}

func (p *PortForwarding) String() string {
	if p.Enabled {
		return "Port forwarding: on, saved in " + p.Filepath
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
func GetPIASettings(envParams libparams.EnvParams) (settings PIA, err error) {
	settings.User, err = params.GetUser(envParams)
	if err != nil {
		return settings, err
	}
	settings.Password, err = params.GetPassword(envParams)
	if err != nil {
		return settings, err
	}
	settings.Encryption, err = params.GetPIAEncryption(envParams)
	if err != nil {
		return settings, err
	}
	settings.Region, err = params.GetPIARegion(envParams)
	if err != nil {
		return settings, err
	}
	settings.PortForwarding.Enabled, err = params.GetPortForwarding(envParams)
	if err != nil {
		return settings, err
	}
	if settings.PortForwarding.Enabled {
		settings.PortForwarding.Filepath, err = params.GetPortForwardingStatusFilepath(envParams)
		if err != nil {
			return settings, err
		}
	}
	return settings, nil
}
