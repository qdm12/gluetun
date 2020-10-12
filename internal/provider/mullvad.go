package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
)

type mullvad struct {
	servers    []models.MullvadServer
	randSource rand.Source
}

func newMullvad(servers []models.MullvadServer, timeNow timeNowFunc) *mullvad {
	return &mullvad{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (m *mullvad) filterServers(country, city, isp string) (servers []models.MullvadServer) {
	for i, server := range m.servers {
		if len(country) == 0 {
			server.Country = ""
		}
		if len(city) == 0 {
			server.City = ""
		}
		if len(isp) == 0 {
			server.ISP = ""
		}
		if strings.EqualFold(server.Country, country) &&
			strings.EqualFold(server.City, city) &&
			strings.EqualFold(server.ISP, isp) {
			servers = append(servers, m.servers[i])
		}
	}
	return servers
}

func (m *mullvad) GetOpenVPNConnection(selection models.ServerSelection) (connection models.OpenVPNConnection, err error) {
	servers := m.filterServers(selection.Country, selection.City, selection.ISP)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for country %q, city %q and ISP %q", selection.Country, selection.City, selection.ISP)
	}

	var defaultPort uint16 = 1194
	if selection.Protocol == constants.TCP {
		defaultPort = 443
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		port := defaultPort
		if selection.CustomPort > 0 {
			port = selection.CustomPort
		}
		for _, IP := range server.IPs {
			if selection.TargetIP != nil {
				if selection.TargetIP.Equal(IP) {
					return models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol}, nil
				}
			} else {
				connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
			}
		}
	}

	if selection.TargetIP != nil {
		return connection, fmt.Errorf("target IP address %q not found in IP addresses", selection.TargetIP)
	}

	return pickRandomConnection(connections, m.randSource), nil
}

func (m *mullvad) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Mullvad specific
		"ping 10",
		"ping-restart 60",
		"sndbuf 524288",
		"rcvbuf 524288",
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA",
		"fast-io",
		"script-security 2",

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		fmt.Sprintf("cipher %s", cipher),
	}
	if extras.OpenVPNIPv6 {
		lines = append(lines, "tun-ipv6")
	} else {
		lines = append(lines, `pull-filter ignore "route-ipv6"`)
		lines = append(lines, `pull-filter ignore "ifconfig-ipv6"`)
	}
	if !root {
		lines = append(lines, "user nonrootuser")
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.MullvadCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
		"",
	}...)
	return lines
}

func (m *mullvad) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for mullvad")
}
