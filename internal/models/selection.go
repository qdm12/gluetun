package models

import (
	"net"
)

type ServerSelection struct {
	// Common
	Protocol NetworkProtocol `json:"network_protocol"`
	TargetIP net.IP          `json:"target_ip,omitempty"`

	// Cyberghost, PIA, Surfshark, Windscribe, Vyprvpn, NordVPN
	Regions []string `json:"regions"`

	// Cyberghost
	Group string `json:"group"`

	Countries []string `json:"countries"` // Mullvad, PureVPN
	Cities    []string `json:"cities"`    // Mullvad, PureVPN, Windscribe
	Hostnames []string `json:"hostnames"` // Windscribe, Privado

	// Mullvad
	ISPs  []string `json:"isps"`
	Owned bool     `json:"owned"`

	// Mullvad, Windscribe, PIA
	CustomPort uint16 `json:"custom_port"`

	// NordVPN
	Numbers []uint16 `json:"numbers"`

	// PIA
	EncryptionPreset string `json:"encryption_preset"`
}

type ExtraConfigOptions struct {
	ClientCertificate string `json:"-"`                 // Cyberghost
	ClientKey         string `json:"-"`                 // Cyberghost
	EncryptionPreset  string `json:"encryption_preset"` // PIA
	OpenVPNIPv6       bool   `json:"openvpn_ipv6"`      // Mullvad
}

// PortForwarding contains settings for port forwarding.
type PortForwarding struct {
	Enabled  bool     `json:"enabled"`
	Filepath Filepath `json:"filepath"`
}
