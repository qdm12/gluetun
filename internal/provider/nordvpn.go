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

func (n *nordvpn) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) { //nolint:dupl
	var IP net.IP
	for _, server := range constants.NordvpnServers() {
		if strings.EqualFold(server.Region, selection.Region) {
			IP = server.IP
			break
		}
	}
	if IP == nil {
		return nil, fmt.Errorf("no IP found for server %q", selection.Region)
	}
	if selection.TargetIP != nil && !selection.TargetIP.Equal(IP) {
		return nil, fmt.Errorf("target IP address %s does not match IP address %s", selection.TargetIP, IP)
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
