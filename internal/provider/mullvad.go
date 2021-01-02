package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
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

func (m *mullvad) filterServers(countries, cities, isps []string, owned bool) (servers []models.MullvadServer) {
	for _, server := range m.servers {
		switch {
		case
			filterByPossibilities(server.Country, countries),
			filterByPossibilities(server.City, cities),
			filterByPossibilities(server.ISP, isps),
			owned && !server.Owned:
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (m *mullvad) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var defaultPort uint16 = 1194
	if selection.Protocol == constants.TCP {
		defaultPort = 443
	}
	port := defaultPort
	if selection.CustomPort > 0 {
		port = selection.CustomPort
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := m.filterServers(selection.Countries, selection.Cities, selection.ISPs, selection.Owned)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for countries %s, cities %s, ISPs %s and owned %t",
			commaJoin(selection.Countries), commaJoin(selection.Cities), commaJoin(selection.ISPs), selection.Owned)
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, m.randSource), nil
}

func (m *mullvad) BuildConf(connection models.OpenVPNConnection,
	verbosity int, username string, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
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
		lines = append(lines, "user "+username)
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
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for mullvad")
}
