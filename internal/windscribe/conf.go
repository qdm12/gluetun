package windscribe

import (
	"fmt"
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetOpenVPNConnections(region models.WindscribeRegion, protocol models.NetworkProtocol, targetIP net.IP) (connections []models.OpenVPNConnection, err error) {
	var subdomain string
	for _, server := range constants.WindscribeServers() {
		if server.Region == region {
			subdomain = server.Subdomain
			break
		}
	}
	if len(subdomain) == 0 {
		return nil, fmt.Errorf("no server found for region %q", region)
	}
	hostname := subdomain + ".windscribe.com"
	IPs, err := c.lookupIP(hostname)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("target IP address %q not found from IP addresses resolved from %s", targetIP, hostname)
		}
		IPs = []net.IP{targetIP}
	}
	var port uint16
	switch protocol {
	case constants.TCP:
		port = 1194
	case constants.UDP:
		port = 443
	default:
		return nil, fmt.Errorf("protocol %q is unknown", protocol)
	}
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: protocol})
	}
	return connections, nil
}

func (c *configurator) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool) (err error) {
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",

		// Windscribe specific
		"reneg-sec 432000",
		"resolv-retry infinite",
		"auth SHA512",
		"cipher AES-256-CBC",
		"keysize 256",
		"comp-lzo",
		"ns-cert-type server",

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
		"",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.WindscribeOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
	}...)
	return c.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}
