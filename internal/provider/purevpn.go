package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type purevpn struct {
	servers    []models.PurevpnServer
	randSource rand.Source
}

func newPurevpn(servers []models.PurevpnServer, timeNow timeNowFunc) *purevpn {
	return &purevpn{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *purevpn) filterServers(regions, countries, cities []string) (servers []models.PurevpnServer) {
	for _, server := range p.servers {
		switch {
		case
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(server.Country, countries),
			filterByPossibilities(server.City, cities):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (p *purevpn) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16
	switch {
	case selection.Protocol == constants.UDP:
		port = 53
	case selection.Protocol == constants.TCP:
		port = 80
	default:
		return connection, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := p.filterServers(selection.Regions, selection.Countries, selection.Cities)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for regions %s, countries %s and cities %s",
			commaJoin(selection.Regions), commaJoin(selection.Countries), commaJoin(selection.Cities))
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, p.randSource), nil
}

func (p *purevpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings settings.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// Purevpn specific
		"key-direction 1",
		"remote-cert-tls server",
		"cipher AES-256-CBC",
		"route-method exe",
		"route-delay 0",
		"route 0.0.0.0 0.0.0.0",
		"script-security 2",

		"ping 10",
		"ping-exit 60",

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
		fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port),
		fmt.Sprintf("cipher %s", settings.Cipher),
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
	if len(settings.Auth) > 0 {
		lines = append(lines, "auth "+settings.Auth)
	}
	if connection.Protocol == constants.UDP {
		lines = append(lines, "explicit-exit-notify")
	}
	return lines
}

func (p *purevpn) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for purevpn")
}
