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

type torguard struct {
	servers    []models.TorguardServer
	randSource rand.Source
}

func newTorguard(servers []models.TorguardServer, timeNow timeNowFunc) *torguard {
	return &torguard{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (t *torguard) filterServers(countries, cities, hostnames []string) (servers []models.TorguardServer) {
	for _, server := range t.servers {
		switch {
		case filterByPossibilities(server.Country, countries):
		case filterByPossibilities(server.City, cities):
		case filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (t *torguard) notFoundErr(selection configuration.ServerSelection) error {
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

func (t *torguard) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16 = 1912
	if selection.CustomPort > 0 {
		port = selection.CustomPort
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := t.filterServers(selection.Countries, selection.Cities, selection.Hostnames)
	if len(servers) == 0 {
		return connection, t.notFoundErr(selection)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connections[i] = models.OpenVPNConnection{
			IP:       servers[i].IP,
			Port:     port,
			Protocol: selection.Protocol,
		}
	}

	return pickRandomConnection(connections, t.randSource), nil
}

func (t *torguard) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256gcm
	}
	if len(settings.Auth) == 0 {
		settings.Auth = sha256
	}

	const defaultMSSFix = 1450
	if settings.MSSFix == 0 {
		settings.MSSFix = defaultMSSFix
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",
		"ping 15",
		"ping-timer-rem",
		"tls-exit",

		// Torguard specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)),
		"reneg-sec 0",
		"fast-io",
		"key-direction 1",
		"script-security 2",
		"ping-restart 0",
		"ncp-disable",
		"compress",
		"keepalive 5 30",
		"sndbuf 393216",
		"rcvbuf 393216",
		// "up /etc/openvpn/update-resolv-conf",
		// "down /etc/openvpn/update-resolv-conf",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"pull-filter ignore \"block-outside-dns\"",
		`pull-filter ignore "ping-restart"`,
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"proto " + connection.Protocol,
		"remote " + connection.IP.String() + " " + strconv.Itoa(int(connection.Port)),
		"cipher " + settings.Cipher,
		"auth " + settings.Auth,
	}

	if !settings.Root {
		lines = append(lines, "user "+username)
	}

	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.TorguardCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)

	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.TorguardOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)

	return lines
}

func (t *torguard) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for torguard")
}
