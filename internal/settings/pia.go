package settings

import (
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
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

func GetPIASettings() (settings PIA, err error) {
	// TODO
	return settings, nil
}
