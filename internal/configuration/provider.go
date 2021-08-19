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
	Name            string          `json:"name"`
	ServerSelection ServerSelection `json:"server_selection"`
	PortForwarding  PortForwarding  `json:"port_forwarding"`
}

func (settings *Provider) lines() (lines []string) {
	if settings.Name == "" { // custom OpenVPN configuration
		return nil
	}

	lines = append(lines, lastIndent+strings.Title(settings.Name)+" settings:")

	for _, line := range settings.ServerSelection.toLines() {
		lines = append(lines, indent+line)
	}

	if settings.PortForwarding.Enabled { // PIA
		lines = append(lines, indent+lastIndent+"Port forwarding:")
		for _, line := range settings.PortForwarding.lines() {
			lines = append(lines, indent+indent+line)
		}
	}

	return lines
}

var (
	ErrInvalidVPNProvider = errors.New("invalid VPN provider")
)

func (settings *Provider) read(r reader, vpnType string) error {
	err := settings.readVPNServiceProvider(r)
	if err != nil {
		return err
	}

	switch settings.Name {
	case constants.Cyberghost:
		err = settings.readCyberghost(r)
	case constants.Fastestvpn:
		err = settings.readFastestvpn(r)
	case constants.HideMyAss:
		err = settings.readHideMyAss(r)
	case constants.Ipvanish:
		err = settings.readIpvanish(r)
	case constants.Ivpn:
		err = settings.readIvpn(r)
	case constants.Mullvad:
		err = settings.readMullvad(r)
	case constants.Nordvpn:
		err = settings.readNordvpn(r)
	case constants.Privado:
		err = settings.readPrivado(r)
	case constants.PrivateInternetAccess:
		err = settings.readPrivateInternetAccess(r)
	case constants.Privatevpn:
		err = settings.readPrivatevpn(r)
	case constants.Protonvpn:
		err = settings.readProtonvpn(r)
	case constants.Purevpn:
		err = settings.readPurevpn(r)
	case constants.Surfshark:
		err = settings.readSurfshark(r)
	case constants.Torguard:
		err = settings.readTorguard(r)
	case constants.VPNUnlimited:
		err = settings.readVPNUnlimited(r)
	case constants.Vyprvpn:
		err = settings.readVyprvpn(r)
	case constants.Windscribe:
		err = settings.readWindscribe(r)
	default:
		return fmt.Errorf("%w: %s", ErrInvalidVPNProvider, settings.Name)
	}
	if err != nil {
		return err
	}

	settings.ServerSelection.VPN = vpnType
	return nil
}

func (settings *Provider) readVPNServiceProvider(r reader) (err error) {
	allowedVPNServiceProviders := []string{
		"cyberghost", "fastestvpn", "hidemyass", "ipvanish", "ivpn", "mullvad", "nordvpn",
		"privado", "pia", "private internet access", "privatevpn", "protonvpn",
		"purevpn", "surfshark", "torguard", constants.VPNUnlimited, "vyprvpn", "windscribe"}

	vpnsp, err := r.env.Inside("VPNSP", allowedVPNServiceProviders,
		params.Default("private internet access"))
	if err != nil {
		return fmt.Errorf("environment variable VPNSP: %w", err)
	}
	if vpnsp == "pia" { // retro compatibility
		vpnsp = "private internet access"
	}
	settings.Name = vpnsp

	return nil
}

func commaJoin(slice []string) string {
	return strings.Join(slice, ", ")
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
