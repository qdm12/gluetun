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

func (n *nordvpn) filterServers(regions []string, protocol models.NetworkProtocol, numbers []uint16) (
	servers []models.NordvpnServer) {
	numbersStr := make([]string, len(numbers))
	for i := range numbers {
		numbersStr[i] = fmt.Sprintf("%d", numbers[i])
	}
	for _, server := range n.servers {
		numberStr := fmt.Sprintf("%d", server.Number)
		switch {
		case
			protocol == constants.TCP && !server.TCP,
			protocol == constants.UDP && !server.UDP,
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(numberStr, numbersStr):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (n *nordvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16
	switch {
	case selection.Protocol == constants.UDP:
		port = 1194
	case selection.Protocol == constants.TCP:
		port = 443
	default:
		return connection, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := n.filterServers(selection.Regions, selection.Protocol, selection.Numbers)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for region %s, protocol %s and numbers %v",
			commaJoin(selection.Regions), selection.Protocol, selection.Numbers)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connections[i] = models.OpenVPNConnection{IP: servers[i].IP, Port: port, Protocol: selection.Protocol}
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
		"ping 10",
		"ping-exit 60",
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

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		`pull-filter ignore "ping-restart"`,
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", settings.Verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port),
		fmt.Sprintf("cipher %s", settings.Cipher),
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
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for nordvpn")
}
