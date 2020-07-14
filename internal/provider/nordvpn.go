package provider

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type nordvpn struct{}

func newNordvpn() *nordvpn {
	return &nordvpn{}
}

func findServers(selection models.ServerSelection) (servers []models.NordvpnServer) {
	for _, server := range constants.NordvpnServers() {
		if strings.EqualFold(server.Region, selection.Region) {
			if (selection.Protocol == constants.TCP && !server.TCP) || (selection.Protocol == constants.UDP && !server.UDP) {
				continue
			}
			if selection.Number > 0 && server.Number == selection.Number {
				return []models.NordvpnServer{server}
			}
			servers = append(servers, server)
		}
	}
	return servers
}

func extractIPsFromServers(servers []models.NordvpnServer) (ips []net.IP) {
	ips = make([]net.IP, len(servers))
	for i := range servers {
		ips[i] = servers[i].IP
	}
	return ips
}

func targetIPInIps(targetIP net.IP, ips []net.IP) error {
	for i := range ips {
		if targetIP.Equal(ips[i]) {
			return nil
		}
	}
	ipsString := make([]string, len(ips))
	for i := range ips {
		ipsString[i] = ips[i].String()
	}
	return fmt.Errorf("target IP address %s not found in IP addresses %s", targetIP, strings.Join(ipsString, ", "))
}

func (n *nordvpn) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) { //nolint:dupl
	servers := findServers(selection)
	ips := extractIPsFromServers(servers)
	if len(ips) == 0 {
		if selection.Number > 0 {
			return nil, fmt.Errorf("no IP found for region %q, protocol %s and number %d", selection.Region, selection.Protocol, selection.Number)
		}
		return nil, fmt.Errorf("no IP found for region %q, protocol %s", selection.Region, selection.Protocol)
	}
	var IP net.IP
	if selection.TargetIP != nil {
		if err := targetIPInIps(selection.TargetIP, ips); err != nil {
			return nil, err
		}
		IP = selection.TargetIP
	} else {
		IP = ips[0]
	}
	var port uint16
	switch {
	case selection.Protocol == constants.UDP:
		port = 1194
	case selection.Protocol == constants.TCP:
		port = 443
	default:
		return nil, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}
	return []models.OpenVPNConnection{{IP: IP, Port: port, Protocol: selection.Protocol}}, nil
}

func (n *nordvpn) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) { //nolint:dupl
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = "sha512"
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Nordvpn specific
		"resolv-retry infinite",
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix 1450",
		"ping 15",
		"ping-restart 0",
		"ping-timer-rem",
		"reneg-sec 0",
		"comp-lzo no",
		"fast-io",
		"key-direction 1",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connections[0].Protocol)),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if !root {
		lines = append(lines, "user nonrootuser")
	}
	for _, connection := range connections {
		lines = append(lines, fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port))
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.NordvpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.NordvpnOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return lines
}

func (n *nordvpn) GetPortForward(client network.Client) (port uint16, err error) {
	panic("port forwarding is not supported for nordvpn")
}
