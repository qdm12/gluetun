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
	"github.com/qdm12/gluetun/internal/settings"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/os"
)

type windscribe struct {
	servers    []models.WindscribeServer
	randSource rand.Source
}

func newWindscribe(servers []models.WindscribeServer, timeNow timeNowFunc) *windscribe {
	return &windscribe{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (w *windscribe) filterServers(regions, cities, hostnames []string) (servers []models.WindscribeServer) {
	for _, server := range w.servers {
		switch {
		case
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(server.City, cities),
			filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

//nolint:lll
func (w *windscribe) GetOpenVPNConnection(selection models.ServerSelection) (connection models.OpenVPNConnection, err error) {
	var port uint16
	switch {
	case selection.CustomPort > 0:
		port = selection.CustomPort
	case selection.Protocol == constants.TCP:
		port = 1194
	case selection.Protocol == constants.UDP:
		port = 443
	default:
		return connection, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := w.filterServers(selection.Regions, selection.Cities, selection.Hostnames)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for region %s", commaJoin(selection.Regions))
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for _, server := range servers {
		connections = append(connections, models.OpenVPNConnection{IP: server.IP, Port: port, Protocol: selection.Protocol})
	}

	return pickRandomConnection(connections, w.randSource), nil
}

func (w *windscribe) BuildConf(connection models.OpenVPNConnection,
	username string, settings settings.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}
	if len(settings.Auth) == 0 {
		settings.Auth = "sha512"
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Windscribe specific
		"comp-lzo",
		"key-direction 1",
		"script-security 2",

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
		fmt.Sprintf("cipher %s", settings.Cipher),
		fmt.Sprintf("auth %s", settings.Auth),
	}
	if strings.HasSuffix(settings.Cipher, "-gcm") {
		lines = append(lines, "ncp-ciphers AES-256-GCM:AES-256-CBC:AES-128-GCM")
	}
	if !settings.Root {
		lines = append(lines, "user "+username)
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.WindscribeCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.WindscribeOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return lines
}

func (w *windscribe) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for windscribe")
}
