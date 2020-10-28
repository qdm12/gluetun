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
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
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

func (c *cyberghost) filterServers(regions []string, group string) (servers []models.CyberghostServer) {
	for _, server := range c.servers {
		switch {
		case len(group) > 0 && !strings.EqualFold(group, server.Group),
			filterByPossibilities(server.Region, regions):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (c *cyberghost) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	const httpsPort = 443
	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: httpsPort, Protocol: selection.Protocol}, nil
	}

	servers := c.filterServers(selection.Regions, selection.Group)
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

func (c *cyberghost) BuildConf(connection models.OpenVPNConnection, verbosity,
	uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = "SHA256"
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"persist-tun",
		"remote-cert-tls server",

		// Cyberghost specific
		// "redirect-gateway def1",
		"ncp-disable",
		"ping 5",
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
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if strings.HasSuffix(cipher, "-gcm") {
		lines = append(lines, "ncp-ciphers AES-256-GCM:AES-256-CBC:AES-128-GCM")
	}
	if !root {
		lines = append(lines, "user nonrootuser")
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
		constants.CyberghostClientCertificate,
		"-----END CERTIFICATE-----",
		"</cert>",
	}...)
	lines = append(lines, []string{
		"<key>",
		"-----BEGIN PRIVATE KEY-----",
		extras.ClientKey,
		"-----END PRIVATE KEY-----",
		"</key>",
		"",
	}...)
	return lines
}

func (c *cyberghost) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for cyberghost")
}
