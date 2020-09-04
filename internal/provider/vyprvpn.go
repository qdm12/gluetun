package provider

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

type vyprvpn struct {
	servers []models.VyprvpnServer
}

func newVyprvpn(servers []models.VyprvpnServer) *vyprvpn {
	return &vyprvpn{
		servers: servers,
	}
}

func (v *vyprvpn) filterServers(region string) (servers []models.VyprvpnServer) {
	if len(region) == 0 {
		return v.servers
	}
	for _, server := range v.servers {
		if strings.EqualFold(server.Region, region) {
			return []models.VyprvpnServer{server}
		}
	}
	return nil
}

func (v *vyprvpn) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	servers := v.filterServers(selection.Region)
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server found for region %q", selection.Region)
	}

	var port uint16
	switch {
	case selection.Protocol == constants.TCP:
		return nil, fmt.Errorf("TCP protocol not supported by this VPN provider")
	case selection.Protocol == constants.UDP:
		port = 443
	default:
		return nil, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	for _, server := range servers {
		for _, IP := range server.IPs {
			if selection.TargetIP != nil {
				if selection.TargetIP.Equal(IP) {
					return []models.OpenVPNConnection{{IP: IP, Port: port, Protocol: selection.Protocol}}, nil
				}
			} else {
				connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
			}
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

func (v *vyprvpn) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = "SHA256"
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Vyprvpn specific
		"comp-lzo",
		"keepalive 10 60",
		// "verify-x509-name lu1.vyprvpn.com name",
		"tls-cipher TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",

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
		fmt.Sprintf("proto %s", connections[0].Protocol),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if !root {
		lines = append(lines, "user nonrootuser")
	}
	for _, connection := range connections {
		lines = append(lines, fmt.Sprintf("remote %s %d", connection.IP, connection.Port))
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.VyprvpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	return lines
}

func (v *vyprvpn) GetPortForward(client network.Client) (port uint16, err error) {
	panic("port forwarding is not supported for vyprvpn")
}
