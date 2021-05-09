package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type privado struct {
	servers    []models.PrivadoServer
	randSource rand.Source
}

func newPrivado(servers []models.PrivadoServer, timeNow timeNowFunc) *privado {
	return &privado{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *privado) filterServers(countries, regions, cities, hostnames []string) (servers []models.PrivadoServer) {
	for _, server := range p.servers {
		switch {
		case filterByPossibilities(server.Country, countries),
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(server.City, cities),
			filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (p *privado) notFoundErr(countries, regions, cities, hostnames []string) error {
	var message string

	if len(countries) > 0 {
		message += " + countries " + commaJoin(countries)
	}

	if len(regions) > 0 {
		message += " + regions " + commaJoin(regions)
	}

	if len(cities) > 0 {
		message += " + cities " + commaJoin(cities)
	}

	if len(hostnames) > 0 {
		message += " + hostnames " + commaJoin(hostnames)
	}

	message = "for " + strings.TrimPrefix(message, " +")

	return fmt.Errorf("%w: %s", errNoServerFound, message)
}

func (p *privado) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16 = 1194
	switch selection.Protocol {
	case constants.UDP:
	default:
		return connection, fmt.Errorf("protocol %q is not supported by Privado", selection.Protocol)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := p.filterServers(selection.Countries, selection.Regions,
		selection.Cities, selection.Hostnames)
	if len(servers) == 0 {
		return connection, p.notFoundErr(selection.Countries,
			selection.Regions, selection.Cities, selection.Hostnames)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connection := models.OpenVPNConnection{
			IP:       servers[i].IP,
			Port:     port,
			Protocol: selection.Protocol,
			Hostname: servers[i].Hostname,
		}
		connections[i] = connection
	}

	return pickRandomConnection(connections, p.randSource), nil
}

func (p *privado) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}
	if len(settings.Auth) == 0 {
		settings.Auth = sha256
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// Privado specific
		"tls-cipher TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-DHE-DSS-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA",
		fmt.Sprintf("verify-x509-name %s name", connection.Hostname),

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", settings.Verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		"data-ciphers-fallback " + settings.Cipher,
		"data-ciphers " + settings.Cipher,
		fmt.Sprintf("auth %s", settings.Auth),
	}
	if !settings.Root {
		lines = append(lines, "user "+username)
	}
	if settings.MSSFix > 0 {
		lines = append(lines, "mssfix "+strconv.Itoa(int(settings.MSSFix)))
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.PrivadoCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	return lines
}

func (p *privado) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for privado")
}
