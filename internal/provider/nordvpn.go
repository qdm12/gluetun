package provider

import (
	"context"
	"errors"
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

type nordvpn struct {
	servers    []models.NordvpnServer
	randSource rand.Source
}

func newNordvpn(servers []models.NordvpnServer, timeNow timeNowFunc) *nordvpn {
	return &nordvpn{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (n *nordvpn) filterServers(regions, hostnames, names []string, numbers []uint16, tcp bool) (
	servers []models.NordvpnServer) {
	numbersStr := make([]string, len(numbers))
	for i := range numbers {
		numbersStr[i] = fmt.Sprintf("%d", numbers[i])
	}
	for _, server := range n.servers {
		numberStr := fmt.Sprintf("%d", server.Number)
		switch {
		case
			tcp && !server.TCP,
			!tcp && !server.UDP,
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(server.Hostname, hostnames),
			filterByPossibilities(server.Name, names),
			filterByPossibilities(numberStr, numbersStr):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

var errNoServerFound = errors.New("no server found")

func (n *nordvpn) notFoundErr(selection configuration.ServerSelection) error {
	message := "for protocol " + tcpBoolToProtocol(selection.TCP)

	if len(selection.Regions) > 0 {
		message += " + regions " + commaJoin(selection.Regions)
	}

	if len(selection.Hostnames) > 0 {
		message += " + hostnames " + commaJoin(selection.Hostnames)
	}

	if len(selection.Names) > 0 {
		message += " + names " + commaJoin(selection.Names)
	}

	if len(selection.Numbers) > 0 {
		numbers := make([]string, len(selection.Numbers))
		for i, n := range selection.Numbers {
			numbers[i] = strconv.Itoa(int(n))
		}
		message += " + numbers " + commaJoin(numbers)
	}

	return fmt.Errorf("%w: %s", errNoServerFound, message)
}

func (n *nordvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16 = 1194
	protocol := constants.UDP
	if selection.TCP {
		port = 443
		protocol = constants.TCP
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: protocol}, nil
	}

	servers := n.filterServers(selection.Regions, selection.Hostnames,
		selection.Names, selection.Numbers, selection.TCP)
	if len(servers) == 0 {
		return connection, n.notFoundErr(selection)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connections[i] = models.OpenVPNConnection{IP: servers[i].IP, Port: port, Protocol: protocol}
	}

	return pickRandomConnection(connections, n.randSource), nil
}

func (n *nordvpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}
	if len(settings.Auth) == 0 {
		settings.Auth = "sha512"
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
		"ping-timer-rem",
		"tls-exit",

		// Nordvpn specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)),
		"reneg-sec 0",
		"comp-lzo no",
		"fast-io",
		"key-direction 1",
		"ping 15",
		"ping-restart 0",

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
		"data-ciphers-fallback " + settings.Cipher,
		"data-ciphers " + settings.Cipher,
		fmt.Sprintf("auth %s", settings.Auth),
	}
	if !settings.Root {
		lines = append(lines, "user "+username)
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.NordvpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.NordvpnOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return lines
}

func (n *nordvpn) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for nordvpn")
}
