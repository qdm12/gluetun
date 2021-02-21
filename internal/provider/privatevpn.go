package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type privatevpn struct {
	servers    []models.PrivatevpnServer
	randSource rand.Source
}

func newPrivatevpn(servers []models.PrivatevpnServer, timeNow timeNowFunc) *privatevpn {
	return &privatevpn{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *privatevpn) filterServers(countries, cities, hostnames []string) (servers []models.PrivatevpnServer) {
	for _, server := range p.servers {
		switch {
		case
			filterByPossibilities(server.Country, countries),
			filterByPossibilities(server.City, cities),
			filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (p *privatevpn) notFoundErr(selection configuration.ServerSelection) error {
	message := "no server found for protocol " + selection.Protocol

	if len(selection.Countries) > 0 {
		message += " + countries " + commaJoin(selection.Countries)
	}

	if len(selection.Cities) > 0 {
		message += " + cities " + commaJoin(selection.Cities)
	}

	if len(selection.Hostnames) > 0 {
		message += " + hostnames " + commaJoin(selection.Hostnames)
	}

	return fmt.Errorf(message)
}

func (p *privatevpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16
	if selection.Protocol == constants.TCP {
		port = 443
	} else {
		port = 1194
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := p.filterServers(selection.Countries, selection.Cities, selection.Hostnames)
	if len(servers) == 0 {
		return connection, p.notFoundErr(selection)
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, ip := range server.IPs {
			connection := models.OpenVPNConnection{
				IP:       ip,
				Port:     port,
				Protocol: selection.Protocol,
			}
			connections = append(connections, connection)
		}
	}

	return pickRandomConnection(connections, p.randSource), nil
}

func (p *privatevpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes128gcm
	}
	if len(settings.Auth) == 0 {
		settings.Auth = sha256
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"tls-exit",

		// Privatevpn specific
		"comp-lzo",
		"tun-ipv6",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"pull-filter ignore \"block-outside-dns\"",
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", settings.Verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		fmt.Sprintf("cipher %s", settings.Cipher),
		fmt.Sprintf("auth %s", settings.Auth),
	}
	if connection.Protocol == constants.UDP {
		lines = append(lines, "key-direction 1")
	}
	if !settings.Root {
		lines = append(lines, "user "+username)
	}
	if settings.MSSFix > 0 {
		line := "mssfix " + strconv.Itoa(int(settings.MSSFix))
		lines = append(lines, line)
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.PrivatevpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-crypt>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.PrivatevpnOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-crypt>",
		"",
	}...)
	return lines
}

func (p *privatevpn) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for privatevpn")
}
