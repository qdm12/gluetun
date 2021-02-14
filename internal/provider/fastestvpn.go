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

type fastestvpn struct {
	servers    []models.FastestvpnServer
	randSource rand.Source
}

func newFastestvpn(servers []models.FastestvpnServer, timeNow timeNowFunc) *fastestvpn {
	return &fastestvpn{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (f *fastestvpn) filterServers(countries, hostnames []string, protocol string) (servers []models.FastestvpnServer) {
	var tcp, udp bool
	if protocol == "tcp" {
		tcp = true
	} else {
		udp = true
	}

	for _, server := range f.servers {
		switch {
		case filterByPossibilities(server.Country, countries):
		case filterByPossibilities(server.Hostname, hostnames):
		case tcp && !server.TCP:
		case udp && !server.UDP:
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (f *fastestvpn) notFoundErr(selection configuration.ServerSelection) error {
	message := "no server found for protocol " + selection.Protocol

	if len(selection.Hostnames) > 0 {
		message += " + hostnames " + commaJoin(selection.Hostnames)
	}

	if len(selection.Countries) > 0 {
		message += " + countries " + commaJoin(selection.Countries)
	}

	return fmt.Errorf(message)
}

func (f *fastestvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16 = 4443

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := f.filterServers(selection.Countries, selection.Hostnames, selection.Protocol)
	if len(servers) == 0 {
		return connection, f.notFoundErr(selection)
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connection := models.OpenVPNConnection{
				IP:       IP,
				Port:     port,
				Protocol: selection.Protocol,
			}
			connections = append(connections, connection)
		}
	}

	return pickRandomConnection(connections, f.randSource), nil
}

func (f *fastestvpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}
	if len(settings.Auth) == 0 {
		settings.Auth = sha256
	}
	if settings.MSSFix == 0 {
		settings.MSSFix = 1450
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"ping 15",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// Fastestvpn specific
		"ping-restart 0",
		"tls-client",
		"tls-cipher  TLS-DHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-DHE-RSA-WITH-AES-256-CBC-SHA:TLS-RSA-WITH-CAMELLIA-256-CBC-SHA:TLS-RSA-WITH-AES-256-CBC-SHA", //nolint:lll
		"comp-lzo",
		"key-direction 1",
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)), // defaults to 1450

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
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
		constants.FastestvpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)

	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.FastestvpnOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)

	return lines
}

func (f *fastestvpn) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for fastestvpn")
}
