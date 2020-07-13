package provider

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/network"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

type windscribe struct{}

func newWindscribe() *windscribe {
	return &windscribe{}
}

func (w *windscribe) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.WindscribeServers() {
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
	case selection.CustomPort > 0:
		port = selection.CustomPort
	case selection.Protocol == constants.TCP:
		port = 1194
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

func (w *windscribe) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
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

		// Windscribe specific
		"comp-lzo",
		"key-direction 1",
		"script-security 2",

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
	if strings.HasSuffix(cipher, "-gcm") {
		lines = append(lines, "ncp-ciphers AES-256-GCM:AES-256-CBC:AES-128-GCM")
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
		constants.WindscribeCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.WindscribeOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return lines
}

func (w *windscribe) GetPortForward(client network.Client) (port uint16, err error) {
	panic("port forwarding is not supported for windscribe")
}
