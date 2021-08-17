package configuration

import (
	"fmt"
	"net"

	"github.com/qdm12/golibs/params"
)

type ServerSelection struct { //nolint:maligned
	// Common
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

	// NordVPN
	Numbers []uint16 `json:"numbers"`

	// ProtonVPN
	FreeOnly bool `json:"free_only"`

	// VPNUnlimited
	StreamOnly bool `json:"stream_only"`

	OpenVPN OpenVPNSelection `json:"openvpn"`
}

type OpenVPNSelection struct {
	TCP        bool   `json:"tcp"`               // UDP if TCP is false
	CustomPort uint16 `json:"custom_port"`       // HideMyAss, Mullvad, PIA, ProtonVPN, Windscribe
	EncPreset  string `json:"encryption_preset"` // PIA - needed to get the port number
}

func (settings *OpenVPNSelection) lines() (lines []string) {
	lines = append(lines, lastIndent+"OpenVPN selection:")

	lines = append(lines, indent+lastIndent+"Protocol: "+protoToString(settings.TCP))

	if settings.CustomPort != 0 {
		lines = append(lines, indent+lastIndent+"Custom port: "+fmt.Sprint(settings.CustomPort))
	}

	if settings.EncPreset != "" {
		lines = append(lines, indent+lastIndent+"PIA encryption preset: "+settings.EncPreset)
	}

	return lines
}

func (settings *OpenVPNSelection) readProtocolOnly(env params.Env) (err error) {
	settings.TCP, err = readProtocol(env)
	return err
}

func (settings *OpenVPNSelection) readProtocolAndPort(env params.Env) (err error) {
	settings.TCP, err = readProtocol(env)
	if err != nil {
		return err
	}

	settings.CustomPort, err = readPortOrZero(env, "PORT")
	if err != nil {
		return fmt.Errorf("environment variable PORT: %w", err)
	}

	return nil
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
