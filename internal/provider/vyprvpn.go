package provider

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type vyprvpn struct{}

func newVyprvpn() *vyprvpn {
	return &vyprvpn{}
}

func (s *vyprvpn) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.VyprvpnServers() {
		if strings.EqualFold(server.Region, selection.Region) {
			IPs = server.IPs
		}
	}
	if len(IPs) == 0 {
		return nil, fmt.Errorf("no IP found for region %q", selection.Region)
	}
	if selection.TargetIP != nil {
		found := false
		for i := range IPs {
			if IPs[i].Equal(selection.TargetIP) {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("target IP address %q not found in IP addresses", selection.TargetIP)
		}
		IPs = []net.IP{selection.TargetIP}
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
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
	}
	return connections, nil
}

func (s *vyprvpn) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
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

func (s *vyprvpn) GetPortForward(client network.Client) (port uint16, err error) {
	panic("port forwarding is not supported for vyprvpn")
}
