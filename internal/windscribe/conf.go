package windscribe

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetOpenVPNConnections(region models.WindscribeRegion, protocol models.NetworkProtocol, customPort uint16, targetIP net.IP) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.WindscribeServers() {
		if strings.EqualFold(string(server.Region), string(region)) {
			IPs = server.IPs
		}
	}
	if len(IPs) == 0 {
		return nil, fmt.Errorf("no IP found for region %q", region)
	}
	if targetIP != nil {
		found := false
		for i := range IPs {
			if IPs[i].Equal(targetIP) {
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("target IP address %q not found in IP addresses", targetIP)
		}
		IPs = []net.IP{targetIP}
	}
	var port uint16
	switch {
	case customPort > 0:
		port = customPort
	case protocol == constants.TCP:
		port = 1194
	case protocol == constants.UDP:
		port = 443
	default:
		return nil, fmt.Errorf("protocol %q is unknown", protocol)
	}
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: protocol})
	}
	return connections, nil
}

func (c *configurator) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string) (err error) {
	if len(cipher) == 0 {
		cipher = "AES-256-CBC"
	}
	if len(auth) == 0 {
		auth = "sha512"
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",

		// Windscribe specific
		"resolv-retry infinite",
		"comp-lzo",
		"remote-cert-tls server",
		"key-direction 1",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connections[0].Protocol)),
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
		lines = append(lines, fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port))
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
	return c.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}
