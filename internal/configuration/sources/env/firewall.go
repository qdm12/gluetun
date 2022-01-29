package env

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

func (r *Reader) readFirewall() (firewall settings.Firewall, err error) {
	vpnInputPortStrings := envToCSV("FIREWALL_VPN_INPUT_PORTS")
	firewall.VPNInputPorts, err = stringsToPorts(vpnInputPortStrings)
	if err != nil {
		return firewall, fmt.Errorf("environment variable FIREWALL_VPN_INPUT_PORTS: %w", err)
	}

	inputPortStrings := envToCSV("FIREWALL_INPUT_PORTS")
	firewall.InputPorts, err = stringsToPorts(inputPortStrings)
	if err != nil {
		return firewall, fmt.Errorf("environment variable FIREWALL_INPUT_PORTS: %w", err)
	}

	outboundSubnetsKey, _ := r.getEnvWithRetro("FIREWALL_OUTBOUND_SUBNETS", "EXTRA_SUBNETS")
	outboundSubnetStrings := envToCSV(outboundSubnetsKey)
	firewall.OutboundSubnets, err = stringsToIPNets(outboundSubnetStrings)
	if err != nil {
		return firewall, fmt.Errorf("environment variable %s: %w", outboundSubnetsKey, err)
	}

	firewall.Enabled, err = envToBoolPtr("FIREWALL")
	if err != nil {
		return firewall, fmt.Errorf("environment variable FIREWALL: %w", err)
	}

	firewall.Debug, err = envToBoolPtr("FIREWALL_DEBUG")
	if err != nil {
		return firewall, fmt.Errorf("environment variable FIREWALL_DEBUG: %w", err)
	}

	return firewall, nil
}

var (
	ErrPortParsing = errors.New("cannot parse port")
	ErrPortValue   = errors.New("port value is not valid")
)

func stringsToPorts(ss []string) (ports []uint16, err error) {
	if len(ss) == 0 {
		return nil, nil
	}
	ports = make([]uint16, len(ss))
	for i, s := range ss {
		port, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s",
				ErrPortParsing, s, err)
		} else if port < 1 || port > 65535 {
			return nil, fmt.Errorf("%w: must be between 1 and 65535: %d",
				ErrPortValue, port)
		}
		ports[i] = uint16(port)
	}
	return ports, nil
}

var (
	ErrIPNetParsing = errors.New("cannot parse IP network")
)

func stringsToIPNets(ss []string) (ipNets []net.IPNet, err error) {
	if len(ss) == 0 {
		return nil, nil
	}
	ipNets = make([]net.IPNet, len(ss))
	for i, s := range ss {
		ip, ipNet, err := net.ParseCIDR(s)
		if err != nil {
			return nil, fmt.Errorf("%w: %s: %s",
				ErrIPNetParsing, s, err)
		}
		ipNet.IP = ip
		ipNets[i] = *ipNet
	}
	return ipNets, nil
}
