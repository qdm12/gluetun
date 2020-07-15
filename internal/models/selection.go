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

type ServerSelection struct { //nolint:maligned
	// Common
	Protocol NetworkProtocol
	TargetIP net.IP

	// Cyberghost, PIA, Surfshark, Windscribe, Vyprvpn, NordVPN
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

	// NordVPN
	Number uint16
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
	customPort := ""
	if p.ServerSelection.CustomPort > 0 {
		customPort = fmt.Sprintf("%d", p.ServerSelection.CustomPort)
	}
	number := ""
	if p.ServerSelection.Number > 0 {
		number = fmt.Sprintf("%d", p.ServerSelection.Number)
	}
	switch strings.ToLower(string(p.Name)) {
	case "private internet access":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
			"Encryption preset: "+p.ExtraConfigOptions.EncryptionPreset,
			"Port forwarding: "+p.PortForwarding.String(),
		)
	case "mullvad":
		settingsList = append(settingsList,
			"Country: "+p.ServerSelection.Country,
			"City: "+p.ServerSelection.City,
			"ISP: "+p.ServerSelection.ISP,
			"Custom port: "+customPort,
		)
	case "windscribe":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
			"Custom port: "+customPort,
		)
	case "surfshark":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
		)
	case "cyberghost":
		settingsList = append(settingsList,
			"ClientKey: [redacted]",
			"Group: "+p.ServerSelection.Group,
			"Region: "+p.ServerSelection.Region,
		)
	case "vyprvpn":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
		)
	case "nordvpn":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
			"Number: "+number,
		)
	default:
		settingsList = append(settingsList,
			"<Missing String method, please implement me!>",
		)
	}
	if p.ServerSelection.TargetIP != nil {
		settingsList = append(settingsList,
			"Target IP address: "+string(p.ServerSelection.TargetIP),
		)
	}
	return strings.Join(settingsList, "\n |--")
}
