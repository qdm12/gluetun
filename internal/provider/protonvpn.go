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

type protonvpn struct {
	servers    []models.ProtonvpnServer
	randSource rand.Source
}

func newProtonvpn(servers []models.ProtonvpnServer, timeNow timeNowFunc) *protonvpn {
	return &protonvpn{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (p *protonvpn) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	port, err := p.getPort(selection)
	if err != nil {
		return connection, err
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{
			IP:       selection.TargetIP,
			Port:     port,
			Protocol: selection.Protocol,
		}, nil
	}

	servers := p.filterServers(selection.Countries, selection.Regions,
		selection.Cities, selection.Names, selection.Hostnames)
	if len(servers) == 0 {
		return connection, p.notFoundErr(selection)
	}

	connections := make([]models.OpenVPNConnection, len(servers))
	for i := range servers {
		connections[i] = models.OpenVPNConnection{
			IP:       servers[i].EntryIP,
			Port:     port,
			Protocol: selection.Protocol,
		}
	}

	return pickRandomConnection(connections, p.randSource), nil
}

func (p *protonvpn) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
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
		"tls-exit",

		// Protonvpn specific
		"tun-mtu 1500",
		"tun-mtu-extra 32",
		"mssfix " + strconv.Itoa(int(settings.MSSFix)),
		"reneg-sec 0",
		"fast-io",
		"key-direction 1",
		"pull",
		"comp-lzo no",

		// Added constant values
		"auth-nocache",
		"mute-replay-warnings",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"pull-filter ignore \"block-outside-dns\"",
		`pull-filter ignore "ping-restart"`,
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"proto " + connection.Protocol,
		"remote " + connection.IP.String() + strconv.Itoa(int(connection.Port)),
		"cipher " + settings.Cipher,
		"auth " + settings.Auth,
	}
	if !settings.Root {
		lines = append(lines, "user "+username)
	}
	lines = append(lines, []string{
		"<ca>",
		"-----BEGIN CERTIFICATE-----",
		constants.ProtonvpnCertificate,
		"-----END CERTIFICATE-----",
		"</ca>",
	}...)
	lines = append(lines, []string{
		"<tls-auth>",
		"-----BEGIN OpenVPN Static key V1-----",
		constants.ProtonvpnOpenvpnStaticKeyV1,
		"-----END OpenVPN Static key V1-----",
		"</tls-auth>",
		"",
	}...)
	return lines
}

func (p *protonvpn) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for protonvpn")
}

func (p *protonvpn) getPort(selection configuration.ServerSelection) (port uint16, err error) {
	if selection.CustomPort == 0 {
		switch selection.Protocol {
		case constants.TCP:
			const defaultTCPPort = 443
			return defaultTCPPort, nil
		case constants.UDP:
			const defaultUDPPort = 1194
			return defaultUDPPort, nil
		}
	}

	port = selection.CustomPort
	switch selection.Protocol {
	case constants.TCP:
		switch port {
		case 443, 5995, 8443: //nolint:gomnd
		default:
			return 0, fmt.Errorf("%w: %d for protocol %s",
				ErrInvalidPort, port, selection.Protocol)
		}
	case constants.UDP:
		switch port {
		case 80, 443, 1194, 4569, 5060: //nolint:gomnd
		default:
			return 0, fmt.Errorf("%w: %d for protocol %s",
				ErrInvalidPort, port, selection.Protocol)
		}
	}

	return port, nil
}

func (p *protonvpn) filterServers(countries, regions, cities, names, hostnames []string) (
	servers []models.ProtonvpnServer) {
	for _, server := range p.servers {
		switch {
		case
			filterByPossibilities(server.Country, countries),
			filterByPossibilities(server.Region, regions),
			filterByPossibilities(server.City, cities),
			filterByPossibilities(server.Name, names),
			filterByPossibilities(server.Hostname, hostnames):
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (p *protonvpn) notFoundErr(selection configuration.ServerSelection) error {
	message := "no server found for protocol " + selection.Protocol

	if len(selection.Countries) > 0 {
		message += " + countries " + commaJoin(selection.Countries)
	}

	if len(selection.Regions) > 0 {
		message += " + regions " + commaJoin(selection.Regions)
	}

	if len(selection.Cities) > 0 {
		message += " + cities " + commaJoin(selection.Cities)
	}

	if len(selection.Names) > 0 {
		message += " + names " + commaJoin(selection.Names)
	}

	if len(selection.Hostnames) > 0 {
		message += " + hostnames " + commaJoin(selection.Hostnames)
	}

	return fmt.Errorf(message)
}
