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

func (s *surfshark) BuildConf(connection models.OpenVPNConnection,
	username string, settings settings.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}
	if len(settings.Auth) == 0 {
		settings.Auth = "SHA512"
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

		// Surfshark specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)),
		"reneg-sec 0",
		"fast-io",
		"key-direction 1",
		"script-security 2",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"pull-filter ignore \"block-outside-dns\"",
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
	if !settings.Root {
		lines = append(lines, "user "+username)
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
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for surfshark")
}
