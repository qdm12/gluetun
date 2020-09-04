package provider

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

type nordvpn struct {
	servers []models.NordvpnServer
}

func newNordvpn(servers []models.NordvpnServer) *nordvpn {
	return &nordvpn{
		servers: servers,
	}
}

func (n *nordvpn) filterServers(region string, protocol models.NetworkProtocol, number uint16) (servers []models.NordvpnServer) {
	for i, server := range n.servers {
		if len(region) == 0 {
			server.Region = ""
		}
		if number == 0 {
			server.Number = 0
		}

		if protocol == constants.TCP && !server.TCP {
			continue
		} else if protocol == constants.UDP && !server.UDP {
			continue
		}
		if strings.EqualFold(server.Region, region) && server.Number == number {
			servers = append(servers, n.servers[i])
		}
	}
	return servers
}

func (n *nordvpn) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) { //nolint:dupl
	servers := n.filterServers(selection.Region, selection.Protocol, selection.Number)
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server found for region %q, protocol %s and number %d", selection.Region, selection.Protocol, selection.Number)
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

	for _, server := range servers {
		if selection.TargetIP != nil {
			if selection.TargetIP.Equal(server.IP) {
				return []models.OpenVPNConnection{{IP: server.IP, Port: port, Protocol: selection.Protocol}}, nil
			}
		} else {
			connections = append(connections, models.OpenVPNConnection{IP: server.IP, Port: port, Protocol: selection.Protocol})
		}
	}

	if selection.TargetIP != nil {
		return nil, fmt.Errorf("target IP %s not found in IP addresses", selection.TargetIP)
	}

	if len(connections) > 64 {
		connections = connections[:64]
	}

	return connections, nil
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
