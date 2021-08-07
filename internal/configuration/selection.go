package configuration

import (
	"net"
)

type ServerSelection struct { //nolint:maligned
	// Common
	TCP      bool   `json:"tcp"` // UDP if TCP is false
	TargetIP net.IP `json:"target_ip,omitempty"`
	// TODO comments
	// Cyberghost, PIA, Protonvpn, Surfshark, Windscribe, Vyprvpn, NordVPN
	Regions []string `json:"regions"`

	// Cyberghost
	Groups []string `json:"groups"`

	// Fastestvpn, HideMyAss, IPVanish, IVPN, Mullvad, PrivateVPN, Protonvpn, PureVPN, VPNUnlimited
	Countries []string `json:"countries"`
	// HideMyAss, IPVanish, IVPN, Mullvad, PrivateVPN, Protonvpn, PureVPN, VPNUnlimited, Windscribe
	Cities []string `json:"cities"`
	// Fastestvpn, HideMyAss, IPVanish, IVPN, PrivateVPN, Windscribe, Privado, Protonvpn, VPNUnlimited
	Hostnames []string `json:"hostnames"`
	Names     []string `json:"names"` // Protonvpn

	// Mullvad
	ISPs  []string `json:"isps"`
	Owned bool     `json:"owned"`

	// Mullvad, Windscribe, PIA
	CustomPort uint16 `json:"custom_port"`

	// NordVPN
	Numbers []uint16 `json:"numbers"`

	// PIA
	EncryptionPreset string `json:"encryption_preset"`

	// ProtonVPN
	FreeOnly bool `json:"free_only"`

	// VPNUnlimited
	StreamOnly bool `json:"stream_only"`
}

type ExtraConfigOptions struct {
	ClientCertificate string `json:"-"`                 // Cyberghost
	ClientKey         string `json:"-"`                 // Cyberghost, VPNUnlimited
	EncryptionPreset  string `json:"encryption_preset"` // PIA
	OpenVPNIPv6       bool   `json:"openvpn_ipv6"`      // Mullvad
}

// PortForwarding contains settings for port forwarding.
type PortForwarding struct {
	Enabled  bool   `json:"enabled"`
	Filepath string `json:"filepath"`
}

func (p *PortForwarding) lines() (lines []string) {
	return []string{
		lastIndent + "File path: " + p.Filepath,
	}
}
