package configuration

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

type ServerSelection struct { //nolint:maligned
	// Common
	VPN      string `json:"vpn"` // note: this is required
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

	// Surfshark
	MultiHopOnly bool `json:"multihop_only"`

	OpenVPN   OpenVPNSelection   `json:"openvpn"`
	Wireguard WireguardSelection `json:"wireguard"`
}

func (selection ServerSelection) toLines() (lines []string) {
	if selection.TargetIP != nil {
		lines = append(lines, lastIndent+"Target IP address: "+selection.TargetIP.String())
	}

	if len(selection.Groups) > 0 {
		lines = append(lines, lastIndent+"Server groups: "+commaJoin(selection.Groups))
	}

	if len(selection.Countries) > 0 {
		lines = append(lines, lastIndent+"Countries: "+commaJoin(selection.Countries))
	}

	if len(selection.Regions) > 0 {
		lines = append(lines, lastIndent+"Regions: "+commaJoin(selection.Regions))
	}

	if len(selection.Cities) > 0 {
		lines = append(lines, lastIndent+"Cities: "+commaJoin(selection.Cities))
	}

	if len(selection.ISPs) > 0 {
		lines = append(lines, lastIndent+"ISPs: "+commaJoin(selection.ISPs))
	}

	if selection.FreeOnly {
		lines = append(lines, lastIndent+"Free servers only")
	}

	if selection.StreamOnly {
		lines = append(lines, lastIndent+"Stream servers only")
	}

	if len(selection.Hostnames) > 0 {
		lines = append(lines, lastIndent+"Hostnames: "+commaJoin(selection.Hostnames))
	}

	if len(selection.Names) > 0 {
		lines = append(lines, lastIndent+"Names: "+commaJoin(selection.Names))
	}

	if len(selection.Numbers) > 0 {
		numbersString := make([]string, len(selection.Numbers))
		for i, numberUint16 := range selection.Numbers {
			numbersString[i] = fmt.Sprint(numberUint16)
		}
		lines = append(lines, lastIndent+"Numbers: "+commaJoin(numbersString))
	}

	if selection.VPN == constants.OpenVPN {
		lines = append(lines, selection.OpenVPN.lines()...)
	} else { // wireguard
		lines = append(lines, selection.Wireguard.lines()...)
	}

	return lines
}

type OpenVPNSelection struct {
	ConfFile   string `json:"conf_file"`         // Custom configuration file path
	TCP        bool   `json:"tcp"`               // UDP if TCP is false
	CustomPort uint16 `json:"custom_port"`       // HideMyAss, Mullvad, PIA, ProtonVPN, Windscribe
	EncPreset  string `json:"encryption_preset"` // PIA - needed to get the port number
}

func (settings *OpenVPNSelection) lines() (lines []string) {
	lines = append(lines, lastIndent+"OpenVPN selection:")

	if settings.ConfFile != "" {
		lines = append(lines, indent+lastIndent+"Custom configuration file: "+settings.ConfFile)
	}

	lines = append(lines, indent+lastIndent+"Protocol: "+protoToString(settings.TCP))

	if settings.CustomPort != 0 {
		lines = append(lines, indent+lastIndent+"Custom port: "+fmt.Sprint(settings.CustomPort))
	}

	if settings.EncPreset != "" {
		lines = append(lines, indent+lastIndent+"PIA encryption preset: "+settings.EncPreset)
	}

	return lines
}

func (settings *OpenVPNSelection) readProtocolOnly(env params.Interface) (err error) {
	settings.TCP, err = readProtocol(env)
	return err
}

func (settings *OpenVPNSelection) readProtocolAndPort(env params.Interface) (err error) {
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

type WireguardSelection struct {
	// EndpointPort is a the server port to use for the VPN server.
	// It is optional for Wireguard VPN providers IVPN, Mullvad
	// and Windscribe, and compulsory for the others
	EndpointPort uint16 `json:"port,omitempty"`
	// PublicKey is the server public key.
	// It is only used with VPN providers generating Wireguard
	// configurations specific to each server and user.
	PublicKey string `json:"publickey,omitempty"`
	// EndpointIP is the server endpoint IP address.
	// It is only used with VPN providers generating Wireguard
	// configurations specific to each server and user.
	EndpointIP net.IP `json:"endpoint_ip,omitempty"`
}

func (settings *WireguardSelection) lines() (lines []string) {
	lines = append(lines, lastIndent+"Wireguard selection:")

	if settings.PublicKey != "" {
		lines = append(lines, indent+lastIndent+"Public key: "+settings.PublicKey)
	}

	if settings.EndpointIP != nil {
		endpoint := settings.EndpointIP.String() + ":" + fmt.Sprint(settings.EndpointPort)
		lines = append(lines, indent+lastIndent+"Server endpoint: "+endpoint)
	} else if settings.EndpointPort != 0 {
		lines = append(lines, indent+lastIndent+"Custom port: "+fmt.Sprint(settings.EndpointPort))
	}

	return lines
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
