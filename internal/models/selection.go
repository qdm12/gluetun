package models

import (
	"fmt"
	"net"
	"strings"
)

// ProviderSettings contains settings specific to a VPN provider.
type ProviderSettings struct {
	Name               VPNProvider        `json:"name"`
	ServerSelection    ServerSelection    `json:"serverSelection"`
	ExtraConfigOptions ExtraConfigOptions `json:"extraConfig"`
	PortForwarding     PortForwarding     `json:"portForwarding"`
}

type ServerSelection struct {
	// Common
	Protocol NetworkProtocol `json:"networkProtocol"`
	TargetIP net.IP          `json:"targetIP,omitempty"`

	// Cyberghost, PIA, Surfshark, Windscribe, Vyprvpn, NordVPN
	Regions []string `json:"regions"`

	// Cyberghost
	Group string `json:"group"`

	Countries []string `json:"countries"` // Mullvad, PureVPN
	Cities    []string `json:"cities"`    // Mullvad, PureVPN, Windscribe
	Hostnames []string `json:"hostnames"` // Windscribe

	// Mullvad
	ISPs  []string `json:"isps"`
	Owned bool     `json:"owned"`

	// Mullvad, Windscribe
	CustomPort uint16 `json:"customPort"`

	// NordVPN
	Numbers []uint16 `json:"numbers"`

	// PIA
	EncryptionPreset string `json:"encryptionPreset"`
}

type ExtraConfigOptions struct {
	ClientKey        string `json:"-"`                // Cyberghost
	EncryptionPreset string `json:"encryptionPreset"` // PIA
	OpenVPNIPv6      bool   `json:"openvpnIPv6"`      // Mullvad
}

// PortForwarding contains settings for port forwarding.
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
	numbers := make([]string, len(p.ServerSelection.Numbers))
	for i, number := range p.ServerSelection.Numbers {
		numbers[i] = fmt.Sprintf("%d", number)
	}
	ipv6 := "off"
	if p.ExtraConfigOptions.OpenVPNIPv6 {
		ipv6 = "on"
	}
	switch strings.ToLower(string(p.Name)) {
	case "private internet access old":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
			"Encryption preset: "+p.ExtraConfigOptions.EncryptionPreset,
			"Port forwarding: "+p.PortForwarding.String(),
		)
	case "private internet access":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
			"Encryption preset: "+p.ExtraConfigOptions.EncryptionPreset,
			"Port forwarding: "+p.PortForwarding.String(),
		)
	case "mullvad":
		settingsList = append(settingsList,
			"Countries: "+commaJoin(p.ServerSelection.Countries),
			"Cities: "+commaJoin(p.ServerSelection.Cities),
			"ISPs: "+commaJoin(p.ServerSelection.ISPs),
			"Custom port: "+customPort,
			"IPv6: "+ipv6,
		)
	case "windscribe":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
			"Custom port: "+customPort,
		)
	case "surfshark":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
		)
	case "cyberghost":
		settingsList = append(settingsList,
			"ClientKey: [redacted]",
			"Group: "+p.ServerSelection.Group,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
		)
	case "vyprvpn":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
		)
	case "nordvpn":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
			"Numbers: "+commaJoin(numbers),
		)
	case "purevpn":
		settingsList = append(settingsList,
			"Regions: "+commaJoin(p.ServerSelection.Regions),
			"Countries: "+commaJoin(p.ServerSelection.Countries),
			"Cities: "+commaJoin(p.ServerSelection.Cities),
		)
	case "privado":
		settingsList = append(settingsList,
			"Cities: "+commaJoin(p.ServerSelection.Cities),
			"Server numbers: "+commaJoin(numbers),
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

func commaJoin(slice []string) string {
	return strings.Join(slice, ", ")
}
