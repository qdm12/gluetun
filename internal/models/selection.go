package models

import (
	"fmt"
	"net"
	"strings"
)

// ProviderSettings contains settings specific to a VPN provider
type ProviderSettings struct {
	Name               VPNProvider        `json:"name"`
	ServerSelection    ServerSelection    `json:"serverSelection"`
	ExtraConfigOptions ExtraConfigOptions `json:"extraConfig"`
	PortForwarding     PortForwarding     `json:"portForwarding"`
}

type ServerSelection struct { //nolint:maligned
	// Common
	Protocol NetworkProtocol `json:"networkProtocol"`
	TargetIP net.IP          `json:"targetIP,omitempty"`

	// Cyberghost, PIA, Surfshark, Windscribe, Vyprvpn, NordVPN
	Region string `json:"region"`

	// Cyberghost
	Group string `json:"group"`

	// Mullvad, PureVPN
	Country string `json:"country"`
	City    string `json:"city"`

	// Mullvad
	ISP   string `json:"isp"`
	Owned bool   `json:"owned"`

	// Mullvad, Windscribe
	CustomPort uint16 `json:"customPort"`

	// NordVPN
	Number uint16 `json:"number"`

	// PIA
	EncryptionPreset string `json:"encryptionPreset"`
}

type ExtraConfigOptions struct {
	ClientKey        string `json:"-"`                // Cyberghost
	EncryptionPreset string `json:"encryptionPreset"` // PIA
}

// PortForwarding contains settings for port forwarding
type PortForwarding struct {
	Enabled  bool     `json:"enabled"`
	Filepath Filepath `json:"filepath"`
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
	case "private internet access old":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
			"Encryption preset: "+p.ExtraConfigOptions.EncryptionPreset,
			"Port forwarding: "+p.PortForwarding.String(),
		)
	case "private internet access":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
			"Encryption preset: "+p.ExtraConfigOptions.EncryptionPreset,
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
	case "purevpn":
		settingsList = append(settingsList,
			"Region: "+p.ServerSelection.Region,
			"Country: "+p.ServerSelection.Country,
			"City: "+p.ServerSelection.City,
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
