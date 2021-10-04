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
	err := settings.readVPNServiceProvider(r, vpnType)
	if err != nil {
		return err
	}

	switch settings.Name {
	case constants.Custom:
		err = settings.readCustom(r, vpnType)
	case constants.Cyberghost:
		err = settings.readCyberghost(r)
	case constants.Expressvpn:
		err = settings.readExpressvpn(r)
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
	case constants.Perfectprivacy:
		err = settings.readPerfectPrivacy(r)
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
	case constants.Wevpn:
		err = settings.readWevpn(r)
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

func (settings *Provider) readVPNServiceProvider(r reader, vpnType string) (err error) {
	var allowedVPNServiceProviders []string
	switch vpnType {
	case constants.OpenVPN:
		allowedVPNServiceProviders = []string{
			constants.Custom,
			"cyberghost", constants.Expressvpn, "fastestvpn", "hidemyass", "ipvanish",
			"ivpn", "mullvad", "nordvpn",
			constants.Perfectprivacy, "privado", "pia", "private internet access", "privatevpn", "protonvpn",
			"purevpn", "surfshark", "torguard", constants.VPNUnlimited, "vyprvpn",
			constants.Wevpn, "windscribe"}
	case constants.Wireguard:
		allowedVPNServiceProviders = []string{
			constants.Custom, constants.Ivpn,
			constants.Mullvad, constants.Windscribe,
		}
	}

	vpnsp, err := r.env.Inside("VPNSP", allowedVPNServiceProviders,
		params.Default("private internet access"))
	if err != nil {
		return fmt.Errorf("environment variable VPNSP: %w", err)
	}
	if vpnsp == "pia" { // retro compatibility
		vpnsp = "private internet access"
	}

	if settings.isOpenVPNCustomConfig(r.env) { // retro compatibility
		vpnsp = constants.Custom
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

func readTargetIP(env params.Interface) (targetIP net.IP, err error) {
	targetIP, err = readIP(env, "OPENVPN_TARGET_IP")
	if err != nil {
		return nil, fmt.Errorf("environment variable OPENVPN_TARGET_IP: %w", err)
	}
	return targetIP, nil
}

type openvpnPortValidation struct {
	allAllowed bool
	tcp        bool
	allowedTCP []uint16
	allowedUDP []uint16
}

func readOpenVPNCustomPort(r reader, validation openvpnPortValidation) (
	port uint16, err error) {
	port, err = readPortOrZero(r.env, "OPENVPN_PORT")
	if err != nil {
		return 0, fmt.Errorf("environment variable OPENVPN_PORT: %w", err)
	} else if port == 0 {
		// Try using old variable name
		port, err = readPortOrZero(r.env, "PORT")
		if err != nil {
			r.onRetroActive("PORT", "OPENVPN_PORT")
			return 0, fmt.Errorf("environment variable PORT: %w", err)
		}
	}

	if port == 0 || validation.allAllowed {
		return port, nil
	}

	if validation.tcp {
		for _, allowedPort := range validation.allowedTCP {
			if port == allowedPort {
				return port, nil
			}
		}
		return 0, fmt.Errorf(
			"environment variable PORT: %w: port %d for TCP protocol, can only be one of %s",
			ErrInvalidPort, port, portsToString(validation.allowedTCP))
	}
	for _, allowedPort := range validation.allowedUDP {
		if port == allowedPort {
			return port, nil
		}
	}
	return 0, fmt.Errorf(
		"environment variable PORT: %w: port %d for UDP protocol, can only be one of %s",
		ErrInvalidPort, port, portsToString(validation.allowedUDP))
}

// note: set allowed to an empty slice to allow all valid ports
func readWireguardCustomPort(env params.Interface, allowed []uint16) (port uint16, err error) {
	port, err = readPortOrZero(env, "WIREGUARD_ENDPOINT_PORT")
	if err != nil {
		return 0, fmt.Errorf("environment variable WIREGUARD_ENDPOINT_PORT: %w", err)
	} else if port == 0 {
		port, _ = readPortOrZero(env, "WIREGUARD_PORT")
		if err == nil {
			return port, nil // 0 or WIREGUARD_PORT value
		}
		return 0, nil // default 0
	}

	if len(allowed) == 0 {
		return port, nil
	}

	for i := range allowed {
		if allowed[i] == port {
			return port, nil
		}
	}

	return 0, fmt.Errorf(
		"environment variable WIREGUARD_PORT: %w: port %d, can only be one of %s",
		ErrInvalidPort, port, portsToString(allowed))
}

func portsToString(ports []uint16) string {
	slice := make([]string, len(ports))
	for i := range ports {
		slice[i] = fmt.Sprint(ports[i])
	}
	return strings.Join(slice, ", ")
}

// isOpenVPNCustomConfig is for retro compatibility to set VPNSP=custom
// if OPENVPN_CUSTOM_CONFIG is set.
func (settings Provider) isOpenVPNCustomConfig(env params.Interface) (ok bool) {
	s, _ := env.Get("VPN_TYPE")
	isOpenVPN := s == constants.OpenVPN
	s, _ = env.Get("OPENVPN_CUSTOM_CONFIG")
	return isOpenVPN && s != ""
}
