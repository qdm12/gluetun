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

type cyberghost struct {
	servers    []models.CyberghostServer
	randSource rand.Source
}

func newCyberghost(servers []models.CyberghostServer, timeNow timeNowFunc) *cyberghost {
	return &cyberghost{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (c *cyberghost) filterServers(regions, hostnames []string, group string) (servers []models.CyberghostServer) {
	for _, server := range c.servers {
		switch {
		case group != "" && !strings.EqualFold(group, server.Group),
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (c *cyberghost) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	const httpsPort = 443
	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: httpsPort, Protocol: selection.Protocol}, nil
	}

	servers := c.filterServers(selection.Regions, selection.Hostnames, selection.Group)
	if len(servers) == 0 {
		return connection,
			fmt.Errorf("no server found for regions %s and group %q", commaJoin(selection.Regions), selection.Group)
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: httpsPort, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, c.randSource), nil
}

func (c *cyberghost) BuildConf(connection models.OpenVPNConnection,
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
		"persist-tun",
		"remote-cert-tls server",
		"ping 10",
		"ping-exit 60",
		"ping-timer-rem",
		"tls-exit",

		// Cyberghost specific
		// "redirect-gateway def1",
		"ncp-disable",
		"explicit-exit-notify 2",
		"script-security 2",
		"route-delay 5",

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
	if strings.HasSuffix(settings.Cipher, "-gcm") {
		lines = append(lines, "ncp-ciphers AES-256-GCM:AES-256-CBC:AES-128-GCM")
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
		constants.CyberghostCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<cert>",
		"-----BEGIN CERTIFICATE-----",
		settings.Provider.ExtraConfigOptions.ClientCertificate,
		"-----END CERTIFICATE-----",
		"</cert>",
	}...)
	lines = append(lines, []string{
		"<key>",
		"-----BEGIN PRIVATE KEY-----",
		settings.Provider.ExtraConfigOptions.ClientKey,
		"-----END PRIVATE KEY-----",
		"</key>",
		"",
	}...)
	return lines
}

func (c *cyberghost) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for cyberghost")
}
