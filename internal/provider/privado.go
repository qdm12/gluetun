package provider

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/firewall"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
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

func (s *privado) filterServers(hostnames []string) (servers []models.PrivadoServer) {
	for _, server := range s.servers {
		switch {
		case filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (s *privado) GetOpenVPNConnection(selection models.ServerSelection) (
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

	servers := s.filterServers(selection.Hostnames)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for cities %s and server numbers %v",
			commaJoin(selection.Cities), selection.Numbers)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connection := models.OpenVPNConnection{
			IP:       servers[i].IP,
			Port:     port,
			Protocol: selection.Protocol,
			Hostname: servers[i].Hostname,
		}
		connections = append(connections, connection)
	}

	return pickRandomConnection(connections, s.randSource), nil
}

func (s *privado) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int, root bool,
	cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = sha256
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",

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
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", connection.Protocol),
		fmt.Sprintf("remote %s %d", connection.IP, connection.Port),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if !root {
		lines = append(lines, "user nonrootuser")
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

func (s *privado) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for privado")
}
