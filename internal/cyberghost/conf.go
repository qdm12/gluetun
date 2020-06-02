package cyberghost

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetOpenVPNConnections(group models.CyberghostGroup, region models.CyberghostRegion, protocol models.NetworkProtocol, targetIP net.IP) (connections []models.OpenVPNConnection, err error) {
	var IPs []net.IP
	for _, server := range constants.CyberghostServers() {
		if strings.EqualFold(string(server.Region), string(region)) && strings.EqualFold(string(server.Group), string(group)) {
			IPs = server.IPs
		}
	}
	if len(IPs) == 0 {
		return nil, fmt.Errorf("no IP found for group %q and region %q", group, region)
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
	for _, IP := range IPs {
		connections = append(connections, models.OpenVPNConnection{IP: IP, Port: 1443, Protocol: protocol})
	}
	return connections, nil
}

func (c *configurator) BuildConf(connections []models.OpenVPNConnection, clientKey string, verbosity, uid, gid int, root bool, cipher, auth string) (err error) {
	if len(cipher) == 0 {
		cipher = "AES-256-CBC"
	}
	if len(auth) == 0 {
		auth = "SHA256"
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",
		"remote-cert-tls server",

		// Cyberghost specific
		"resolv-retry infinite",
		"redirect-gateway def1",
		"ncp-disable",
		"ping 5",
		"ping-exit 60",
		"ping-timer-rem",
		"script-security 2",
		"route-delay 5",

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
		constants.CyberghostCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<crt>",
		"-----BEGIN CERTIFICATE-----",
		constants.CyberghostClientCertificate,
		"-----END CERTIFICATE-----",
		"</crt>",
	}...)
	lines = append(lines, []string{
		"<key>",
		"-----BEGIN PRIVATE KEY-----",
		clientKey,
		"-----END PRIVATE KEY-----",
		"</key>",
		"",
	}...)
	return c.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}
