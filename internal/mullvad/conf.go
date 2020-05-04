package mullvad

import (
	"fmt"
	"net"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) GetOpenVPNConnections(country models.MullvadCountry, city models.MullvadCity, provider models.MullvadProvider, protocol models.NetworkProtocol, customPort uint16, targetIP net.IP) (connections []models.OpenVPNConnection, err error) {
	servers := constants.MullvadServerFilter(country, city, provider)
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server found for country %q, city %q and ISP %q", country, city, provider)
	}
	for _, server := range servers {
		port := server.DefaultPort
		if customPort > 0 {
			port = customPort
		}
		for _, IP := range server.IPs {
			if targetIP != nil {
				if targetIP.Equal(IP) {
					return []models.OpenVPNConnection{{IP: IP, Port: port, Protocol: protocol}}, nil
				}
			} else {
				connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: protocol})
			}
		}
	}
	if targetIP != nil {
		return nil, fmt.Errorf("target IP address %q not found in IP addresses", targetIP)
	}
	return connections, nil
}

func (c *configurator) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher string) (err error) {
	if len(connections) == 0 {
		return fmt.Errorf("at least one connection string is expected")
	}
	if len(cipher) == 0 {
		cipher = "AES-256-CBC"
	}
	lines := []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",
		"tls-client",
		"remote-cert-tls server",
		"ping 10",
		"ping-restart 60",

		// Mullvad specific
		"sndbuf 524288",
		"rcvbuf 524288",
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",
		"tun-ipv6",
		"fast-io",

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connections[0].Protocol)),
		fmt.Sprintf("cipher %s", cipher),
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
		constants.MullvadCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
		"",
	}...)
	return c.fileManager.WriteLinesToFile(string(constants.OpenVPNConf), lines, files.Ownership(uid, gid), files.Permissions(0400))
}
