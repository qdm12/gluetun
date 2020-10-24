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

type vyprvpn struct {
	servers    []models.VyprvpnServer
	randSource rand.Source
}

func newVyprvpn(servers []models.VyprvpnServer, timeNow timeNowFunc) *vyprvpn {
	return &vyprvpn{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (v *vyprvpn) filterServers(regions []string) (servers []models.VyprvpnServer) {
	for _, server := range v.servers {
		switch {
		case
			filterByPossibilities(server.Region, regions):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (v *vyprvpn) GetOpenVPNConnection(selection models.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var port uint16
	switch {
	case selection.Protocol == constants.TCP:
		return connection, fmt.Errorf("TCP protocol not supported by this VPN provider")
	case selection.Protocol == constants.UDP:
		port = 443
	default:
		return connection, fmt.Errorf("protocol %q is unknown", selection.Protocol)
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := v.filterServers(selection.Regions)
	if len(servers) == 0 {
		return connection, fmt.Errorf("no server found for region %s", commaJoin(selection.Regions))
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, v.randSource), nil
}

func (v *vyprvpn) BuildConf(connection models.OpenVPNConnection, verbosity, uid, gid int,
	root bool, cipher, auth string, extras models.ExtraConfigOptions) (lines []string) {
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
		"remote-cert-tls server",

		// Vyprvpn specific
		"comp-lzo",
		"keepalive 10 60",
		// "verify-x509-name lu1.vyprvpn.com name",
		"tls-cipher TLS-ECDHE-RSA-WITH-AES-256-GCM-SHA384:TLS-DHE-RSA-WITH-AES-256-CBC-SHA256:TLS-DHE-RSA-WITH-AES-256-CBC-SHA", //nolint:lll

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
		constants.VyprvpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	return lines
}

func (v *vyprvpn) PortForward(ctx context.Context, client *http.Client,
	fileManager files.FileManager, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath models.Filepath)) {
	panic("port forwarding is not supported for vyprvpn")
}
