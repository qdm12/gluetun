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

func (n *nordvpn) filterServers(region string, protocol models.NetworkProtocol, number uint16) (servers []models.NordvpnServer) {
	for i, server := range n.servers {
		if len(region) == 0 {
			server.Region = ""
		}
		if number == 0 {
			server.Number = 0
		}

		if protocol == constants.TCP && !server.TCP {
			continue
		} else if protocol == constants.UDP && !server.UDP {
			continue
		}
		if strings.EqualFold(server.Region, region) && server.Number == number {
			servers = append(servers, n.servers[i])
		}
	}
	return servers
}

func (n *nordvpn) GetOpenVPNConnection(selection models.ServerSelection) (connection models.OpenVPNConnection, err error) { //nolint:dupl
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

	servers := n.filterServers(selection.Region, selection.Protocol, selection.Number)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for region %q, protocol %s and number %d", selection.Region, selection.Protocol, selection.Number)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connections = append(connections, models.OpenVPNConnection{IP: servers[i].IP, Port: port, Protocol: selection.Protocol})
	}

	return pickRandomConnection(connections, n.randSource), nil
}

func (n *nordvpn) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int, root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) { //nolint:dupl
	if len(cipher) == 0 {
		cipher = aes256cbc
	}
	if len(auth) == 0 {
		auth = "sha512"
	}
	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"remote-cert-tls server",

		// Nordvpn specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix 1450",
		"ping 15",
		"ping-restart 0",
		"ping-timer-rem",
		"reneg-sec 0",
		"comp-lzo no",
		"fast-io",
		"key-direction 1",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"remote-random",
		"suppress-timestamps",

		// Modified variables
		fmt.Sprintf("verb %d", verbosity),
		fmt.Sprintf("auth-user-pass %s", constants.OpenVPNAuthConf),
		fmt.Sprintf("proto %s", string(connection.Protocol)),
		fmt.Sprintf("remote %s %d", connection.IP.String(), connection.Port),
		fmt.Sprintf("cipher %s", cipher),
		fmt.Sprintf("auth %s", auth),
	}
	if !root {
		lines = append(lines, "user nonrootuser")
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
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for nordvpn")
}
