package configuration

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/params"
)

// Provider contains settings specific to a VPN provider.
type Provider struct {
	Name               models.VPNProvider `json:"name"`
	ServerSelection    ServerSelection    `json:"server_selection"`
	ExtraConfigOptions ExtraConfigOptions `json:"extra_config"`
	PortForwarding     PortForwarding     `json:"port_forwarding"`
}

func (p *Provider) lines() (lines []string) {
	lines = append(lines, lastIndent+strings.Title(string(p.Name))+" settings:")

	lines = append(lines, indent+lastIndent+"Network protocol: "+string(p.ServerSelection.Protocol))

	if p.ServerSelection.TargetIP != nil {
		lines = append(lines, indent+lastIndent+"Target IP address: "+p.ServerSelection.TargetIP.String())
	}

	var providerLines []string
	switch strings.ToLower(string(p.Name)) {
	case "cyberghost":
		providerLines = p.cyberghostLines()
	case "mullvad":
		providerLines = p.mullvadLines()
	case "private internet access":
		providerLines = p.privateinternetaccessLines()
	case "windscribe":
		providerLines = p.windscribeLines()
	case "surfshark":
		providerLines = p.surfsharkLines()
	case "vyprvpn":
		providerLines = p.vyprvpnLines()
	case "nordvpn":
		providerLines = p.nordvpnLines()
	case "purevpn":
		providerLines = p.purevpnLines()
	case "privado":
		providerLines = p.privadoLines()
	default:
		panic("Missing lines method for provider " + p.Name + "! Please create a Github issue.")
	}

	for _, line := range providerLines {
		lines = append(lines, indent+line)
	}

	return lines
}

func commaJoin(slice []string) string {
	return strings.Join(slice, ", ")
}

func readProtocol(env params.Env) (protocol models.NetworkProtocol, err error) {
	s, err := env.Inside("PROTOCOL",
		[]string{string(constants.TCP), string(constants.UDP)},
		params.Default(string(constants.UDP)))
	if err != nil {
		return "", err
	}
	return models.NetworkProtocol(s), nil
}

func readTargetIP(env params.Env) (targetIP net.IP, err error) {
	return readIP(env, "OPENVPN_TARGET_IP")
}

var (
	ErrInvalidProtocol = errors.New("invalid network protocol")
)

func readCustomPort(env params.Env, protocol models.NetworkProtocol,
	allowedTCP, allowedUDP []uint16) (port uint16, err error) {
	port, err = env.Port("PORT", params.Default("0"))
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
			if allowedTCP[i] == port {
				return port, nil
			}
		}
		return 0, fmt.Errorf("%w: port %d for UDP protocol", ErrInvalidPort, port)
	default:
		return 0, fmt.Errorf("%w: %s", ErrInvalidProtocol, protocol)
	}
}
