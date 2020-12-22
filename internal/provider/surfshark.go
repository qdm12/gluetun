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

type surfshark struct {
	servers    []models.SurfsharkServer
	randSource rand.Source
}

func newSurfshark(servers []models.SurfsharkServer, timeNow timeNowFunc) *surfshark {
	return &surfshark{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (s *surfshark) filterServers(regions []string) (servers []models.SurfsharkServer) {
	for _, server := range s.servers {
		switch {
		case
			filterByPossibilities(server.Region, regions):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (s *surfshark) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16
	switch {
	case selection.Protocol == constants.TCP:
		port = 1443
	case selection.Protocol == constants.UDP:
		port = 1194
	default:
		return connection, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := s.filterServers(selection.Regions)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for region %s", commaJoin(selection.Regions))
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	if selection.TargetIP != nil {
		return connection, fmt.Errorf("target IP %s not found in IP addresses", selection.TargetIP)
	}

	return pickRandomConnection(connections, s.randSource), nil
}

func (s *surfshark) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int, root bool,
	cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = "SHA512"
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Surfshark specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix 1450",
		"ping 15",
		"ping-restart 60",
		"ping-timer-rem",
		"reneg-sec 0",
		"fast-io",
		"key-direction 1",
		"script-security 2",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"pull-filter ignore \"register-dns\"",
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
		constants.SurfsharkCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.SurfsharkOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return lines
}

func (s *surfshark) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for surfshark")
}
