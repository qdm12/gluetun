package provider

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

type purevpn struct {
	servers []models.PurevpnServer
}

func newPurevpn(servers []models.PurevpnServer) *purevpn {
	return &purevpn{
		servers: servers,
	}
}

func (p *purevpn) filterServers(region, country, city string) (servers []models.PurevpnServer) {
	for i, server := range p.servers {
		if len(region) == 0 {
			server.Region = ""
		}
		if len(country) == 0 {
			server.Country = ""
		}
		if len(city) == 0 {
			server.City = ""
		}
		if strings.EqualFold(server.Region, region) &&
			strings.EqualFold(server.Country, country) &&
			strings.EqualFold(server.City, city) {
			servers = append(servers, p.servers[i])
		}
	}
	return servers
}

func (p *purevpn) GetOpenVPNConnections(selection models.ServerSelection) (connections []models.OpenVPNConnection, err error) { //nolint:dupl
	servers := p.filterServers(selection.Region, selection.Country, selection.City)
	if len(servers) == 0 {
		return nil, fmt.Errorf("no server found for region %q, country %q and city %q", selection.Region, selection.Country, selection.City)
	}

	var port uint16
	switch {
	case selection.Protocol == constants.UDP:
		port = 53
	case selection.Protocol == constants.TCP:
		port = 80
	default:
		return nil, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	for _, server := range servers {
		for _, IP := range server.IPs {
			if selection.TargetIP != nil {
				if IP.Equal(selection.TargetIP) {
					return []models.OpenVPNConnection{{IP: IP, Port: port, Protocol: selection.Protocol}}, nil
				}
			} else {
				connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
			}
		}
	}

	if selection.TargetIP != nil {
		return nil, fmt.Errorf("target IP address %q not found in IP addresses", selection.TargetIP)
	}

	if len(connections) > 64 {
		connections = connections[:64]
	}

	return connections, nil
}

func (p *purevpn) BuildConf(connections []models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) { //nolint:dupl
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Purevpn specific
		"key-direction 1",
		"remote-cert-tls server",
		"cipher AES-256-CBC",
		"route-method exe",
		"route-delay 0",
		"route 0.0.0.0 0.0.0.0",
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
		constants.PurevpnCertificateAuthority,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<cert>",
		"-----BEGIN CERTIFICATE-----",
		constants.PurevpnCertificate,
		"-----END CERTIFICATE-----",
		"</cert>",
	}...)
	lines = append(lines, []string{
		"<key>",
		"-----BEGIN PRIVATE KEY-----",
		constants.PurevpnKey,
		"-----END PRIVATE KEY-----",
		"</key>",
		"",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.PurevpnOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	if len(auth) > 0 {
		lines = append(lines, "auth "+auth)
	}
	if connections[0].Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}
	return lines
}

func (p *purevpn) GetPortForward(client network.Client) (port uint16, err error) {
	panic("port forwarding is not supported for purevpn")
}
