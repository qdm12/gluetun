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

	selection := settings.ServerSelection

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

var (
	ErrInvalidVPNProvider = errors.New("invalid VPN provider")
)

func (settings *Provider) read(r reader) error {
	err := settings.readVPNServiceProvider(r)
	if err != nil {
		return err
	}

	var readProvider func(r reader) error
	switch settings.Name {
	case "": // custom config
		readProvider = func(r reader) error { return nil }
	case constants.Cyberghost:
		readProvider = settings.readCyberghost
	case constants.Fastestvpn:
		readProvider = settings.readFastestvpn
	case constants.HideMyAss:
		readProvider = settings.readHideMyAss
	case constants.Ipvanish:
		readProvider = settings.readIpvanish
	case constants.Ivpn:
		readProvider = settings.readIvpn
	case constants.Mullvad:
		readProvider = settings.readMullvad
	case constants.Nordvpn:
		readProvider = settings.readNordvpn
	case constants.Privado:
		readProvider = settings.readPrivado
	case constants.PrivateInternetAccess:
		readProvider = settings.readPrivateInternetAccess
	case constants.Privatevpn:
		readProvider = settings.readPrivatevpn
	case constants.Protonvpn:
		readProvider = settings.readProtonvpn
	case constants.Purevpn:
		readProvider = settings.readPurevpn
	case constants.Surfshark:
		readProvider = settings.readSurfshark
	case constants.Torguard:
		readProvider = settings.readTorguard
	case constants.VPNUnlimited:
		readProvider = settings.readVPNUnlimited
	case constants.Vyprvpn:
		readProvider = settings.readVyprvpn
	case constants.Windscribe:
		readProvider = settings.readWindscribe
	default:
		return fmt.Errorf("%w: %s", ErrInvalidVPNProvider, settings.Name)
	}
	return readProvider(r)
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
