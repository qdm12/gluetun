package configuration

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/params"
)

// Provider contains settings specific to a VPN provider.
type Provider struct {
	Name               string             `json:"name"`
	ServerSelection    ServerSelection    `json:"server_selection"`
	ExtraConfigOptions ExtraConfigOptions `json:"extra_config"`
	PortForwarding     PortForwarding     `json:"port_forwarding"`
}

func (settings *Provider) lines() (lines []string) {
	lines = append(lines, lastIndent+strings.Title(settings.Name)+" settings:")

	selection := settings.ServerSelection

	lines = append(lines, indent+lastIndent+"Network protocol: "+protoToString(selection.TCP))

	if selection.TargetIP != nil {
		lines = append(lines, indent+lastIndent+"Target IP address: "+selection.TargetIP.String())
	}

	var providerLines []string
	switch strings.ToLower(settings.Name) {
	case "cyberghost":
		providerLines = settings.cyberghostLines()
	case "fastestvpn":
		providerLines = settings.fastestvpnLines()
	case "hidemyass":
		providerLines = settings.hideMyAssLines()
	case "ipvanish":
		providerLines = settings.ipvanishLines()
	case "ivpn":
		providerLines = settings.ivpnLines()
	case "mullvad":
		providerLines = settings.mullvadLines()
	case "nordvpn":
		providerLines = settings.nordvpnLines()
	case "privado":
		providerLines = settings.privadoLines()
	case "privatevpn":
		providerLines = settings.privatevpnLines()
	case "private internet access":
		providerLines = settings.privateinternetaccessLines()
	case "protonvpn":
		providerLines = settings.protonvpnLines()
	case "purevpn":
		providerLines = settings.purevpnLines()
	case "surfshark":
		providerLines = settings.surfsharkLines()
	case "torguard":
		providerLines = settings.torguardLines()
	case strings.ToLower(constants.VPNUnlimited):
		providerLines = settings.vpnUnlimitedLines()
	case "vyprvpn":
		providerLines = settings.vyprvpnLines()
	case "windscribe":
		providerLines = settings.windscribeLines()
	default:
		panic(`Missing lines method for provider "` +
			settings.Name + `"! Please create a Github issue.`)
	}

	for _, line := range providerLines {
		lines = append(lines, indent+line)
	}

	return lines
}

func commaJoin(slice []string) string {
	return strings.Join(slice, ", ")
}

func readProtocol(env params.Env) (tcp bool, err error) {
	protocol, err := env.Inside("PROTOCOL", []string{constants.TCP, constants.UDP}, params.Default(constants.UDP))
	if err != nil {
		return false, fmt.Errorf("environment variable PROTOCOL: %w", err)
	}
	return protocol == constants.TCP, nil
}

func protoToString(tcp bool) string {
	if tcp {
		return constants.TCP
	}
	return constants.UDP
}

func readTargetIP(env params.Env) (targetIP net.IP, err error) {
	targetIP, err = readIP(env, "OPENVPN_TARGET_IP")
	if err != nil {
		return nil, fmt.Errorf("environment variable OPENVPN_TARGET_IP: %w", err)
	}
	return targetIP, nil
}

func readCustomPort(env params.Env, tcp bool,
	allowedTCP, allowedUDP []uint16) (port uint16, err error) {
	port, err = readPortOrZero(env, "PORT")
	if err != nil {
		return 0, fmt.Errorf("environment variable PORT: %w", err)
	} else if port == 0 {
		return 0, nil
	}

	if tcp {
		for i := range allowedTCP {
			if allowedTCP[i] == port {
				return port, nil
			}
		}
		return 0, fmt.Errorf("environment variable PORT: %w: port %d for TCP protocol", ErrInvalidPort, port)
	}
	for i := range allowedUDP {
		if allowedUDP[i] == port {
			return port, nil
		}
	}
	return 0, fmt.Errorf("environment variable PORT: %w: port %d for UDP protocol", ErrInvalidPort, port)
}
