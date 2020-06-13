package models

import (
	"fmt"
	"net"
	"strings"
)

// ProviderSettings contains settings specific to a VPN provider
type ProviderSettings struct {
	Name               VPNProvider
	ServerSelection    ServerSelection
	ExtraConfigOptions ExtraConfigOptions
	PortForwarding     PortForwarding
}

type ServerSelection struct {
	// Common
	Protocol NetworkProtocol
	TargetIP net.IP

	// Cyberghost, PIA, Surfshark, Windscribe
	Region string

	// Cyberghost
	Group string

	// Mullvad
	Country string
	City    string
	ISP     string
	Owned   bool

	// Mullvad, Windscribe
	CustomPort uint16

	// PIA
	EncryptionPreset string
}

type ExtraConfigOptions struct {
	ClientKey        string // Cyberghost
	EncryptionPreset string // PIA
}

// PortForwarding contains settings for port forwarding
type PortForwarding struct {
	Enabled  bool
	Filepath Filepath
}

func (p *PortForwarding) String() string {
	if p.Enabled {
		return fmt.Sprintf("on, saved in %s", p.Filepath)
	}
	return "off"
}

func (p *ProviderSettings) String() string {
	settingsList := []string{
		fmt.Sprintf("%s settings:", strings.Title(string(p.Name))),
		"Network protocol: " + string(p.ServerSelection.Protocol),
	}
	switch strings.ToLower(string(p.Name)) {
	case "private internet access":
		settingsList = []string{
			"Region: " + p.ServerSelection.Region,
			"Encryption preset: " + p.ExtraConfigOptions.EncryptionPreset,
			"Port forwarding: " + p.PortForwarding.String(),
		}
	case "mullvad":
		settingsList = []string{
			"Country: " + p.ServerSelection.Country,
			"City: " + p.ServerSelection.City,
			"ISP: " + p.ServerSelection.ISP,
			"Custom port: " + string(p.ServerSelection.CustomPort),
		}
	case "windscribe":
		settingsList = []string{
			"Region: " + p.ServerSelection.Region,
			"Custom port: " + string(p.ServerSelection.CustomPort),
		}
	case "surfshark":
		settingsList = []string{
			"Region: " + p.ServerSelection.Region,
		}
	case "cyberghost":
		settingsList = []string{
			"ClientKey: [redacted]",
			"Group: " + p.ServerSelection.Group,
			"Region: " + p.ServerSelection.Region,
		}
	}
	if p.ServerSelection.TargetIP != nil {
		settingsList = append(settingsList,
			"Target IP address: "+string(p.ServerSelection.TargetIP),
		)
	}
	return strings.Join(settingsList, "\n |--")
}
