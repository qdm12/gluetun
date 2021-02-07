package configuration

import (
	"errors"
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

	lines = append(lines, indent+lastIndent+"Network protocol: "+settings.ServerSelection.Protocol)

	if settings.ServerSelection.TargetIP != nil {
		lines = append(lines, indent+lastIndent+"Target IP address: "+settings.ServerSelection.TargetIP.String())
	}

	var providerLines []string
	switch strings.ToLower(settings.Name) {
	case "cyberghost":
		providerLines = settings.cyberghostLines()
	case "mullvad":
		providerLines = settings.mullvadLines()
	case "nordvpn":
		providerLines = settings.nordvpnLines()
	case "privado":
		providerLines = settings.privadoLines()
	case "private internet access":
		providerLines = settings.privateinternetaccessLines()
	case "purevpn":
		providerLines = settings.purevpnLines()
	case "surfshark":
		providerLines = settings.surfsharkLines()
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

func readProtocol(env params.Env) (protocol string, err error) {
	return env.Inside("PROTOCOL", []string{constants.TCP, constants.UDP}, params.Default(constants.UDP))
}

func readTargetIP(env params.Env) (targetIP net.IP, err error) {
	return readIP(env, "OPENVPN_TARGET_IP")
}

var (
	ErrInvalidProtocol = errors.New("invalid network protocol")
)

func readCustomPort(env params.Env, protocol string,
	allowedTCP, allowedUDP []uint16) (port uint16, err error) {
	port, err = readPortOrZero(env, "PORT")
	if err != nil {
		return 0, err
	} else if port == 0 {
		return 0, nil
	}

	switch protocol {
	case constants.TCP:
		for i := range allowedTCP {
			if allowedTCP[i] == port {
				return port, nil
			}
		}
		return 0, fmt.Errorf("%w: port %d for TCP protocol", ErrInvalidPort, port)
	case constants.UDP:
		for i := range allowedUDP {
			if allowedUDP[i] == port {
				return port, nil
			}
		}
		return 0, fmt.Errorf("%w: port %d for UDP protocol", ErrInvalidPort, port)
	default:
		return 0, fmt.Errorf("%w: %s", ErrInvalidProtocol, protocol)
	}
}
