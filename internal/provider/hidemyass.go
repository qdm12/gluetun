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

type hideMyAss struct {
	servers    []models.HideMyAssServer
	randSource rand.Source
}

func newHideMyAss(servers []models.HideMyAssServer, timeNow timeNowFunc) *hideMyAss {
	return &hideMyAss{
		servers:    servers,
		randSource: rand.NewSource(timeNow().UnixNano()),
	}
}

func (h *hideMyAss) filterServers(countries, cities, hostnames []string,
	protocol string) (servers []models.HideMyAssServer) {
	for _, server := range h.servers {
		switch {
		case
			filterByPossibilities(server.Country, countries),
			filterByPossibilities(server.City, cities),
			filterByPossibilities(server.Hostname, hostnames),
			protocol == constants.TCP && !server.TCP,
			protocol == constants.UDP && !server.UDP:
		default:
			servers = append(servers, server)
		}
	}
	return servers
}

func (h *hideMyAss) notFoundErr(selection configuration.ServerSelection) error {
	var filters []string

	if len(selection.Countries) > 0 {
		filters = append(filters, "countries "+commaJoin(selection.Countries))
	}

	if len(selection.Cities) > 0 {
		filters = append(filters, "countries "+commaJoin(selection.Cities))
	}

	if len(selection.Hostnames) > 0 {
		filters = append(filters, "countries "+commaJoin(selection.Hostnames))
	}

	return fmt.Errorf("%w for %s", ErrNoServerFound, strings.Join(filters, " + "))
}

func (h *hideMyAss) GetOpenVPNConnection(selection configuration.ServerSelection) (
	connection models.OpenVPNConnection, err error) {
	var defaultPort uint16 = 553
	if selection.Protocol == constants.TCP {
		defaultPort = 8080
	}
	port := defaultPort
	if selection.CustomPort > 0 {
		port = selection.CustomPort
	}

	if selection.TargetIP != nil {
		return models.OpenVPNConnection{IP: selection.TargetIP, Port: port, Protocol: selection.Protocol}, nil
	}

	servers := h.filterServers(selection.Countries, selection.Cities, selection.Hostnames, selection.Protocol)
	if len(servers) == 0 {
		return models.OpenVPNConnection{}, h.notFoundErr(selection)
	}

	var connections []models.OpenVPNConnection
	for _, server := range servers {
		for _, IP := range server.IPs {
			connections = append(connections, models.OpenVPNConnection{IP: IP, Port: port, Protocol: selection.Protocol})
		}
	}

	return pickRandomConnection(connections, h.randSource), nil
}

func (h *hideMyAss) BuildConf(connection models.OpenVPNConnection,
	username string, settings configuration.OpenVPN) (lines []string) {
	if len(settings.Cipher) == 0 {
		settings.Cipher = aes256cbc
	}

	lines = []string{
		"client",
		"dev tun",
		"nobind",
		"persist-key",
		"ping 5",
		"ping-exit 30",
		"ping-timer-rem",
		"tls-exit",

		// HideMyAss specific
		"remote-cert-tls server", // updated name of ns-cert-type
		// "route-metric 1",
		"comp-lzo yes",
		"comp-noadapt",

		// Added constant values
		"mute-replay-warnings",
		"auth-nocache",
		"pull-filter ignore \"auth-token\"", // prevent auth failed loops
		"auth-retry nointeract",
		"suppress-timestamps",

		// Modified variables
		"verb " + strconv.Itoa(settings.Verbosity),
		"auth-user-pass " + constants.OpenVPNAuthConf,
		"proto " + connection.Protocol,
		"remote " + connection.IP.String() + strconv.Itoa(int(connection.Port)),
		"data-ciphers-fallback " + settings.Cipher,
		"data-ciphers " + settings.Cipher,
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
		constants.HideMyAssCA,
		"-----END CERTIFICATE-----",
		"</ca>",
		"<cert>",
		"-----BEGIN CERTIFICATE-----",
		constants.HideMyAssCertificate,
		"-----END CERTIFICATE-----",
		"</cert>",
		"<key>",
		"-----BEGIN RSA PRIVATE KEY-----",
		constants.HideMyAssRSAPrivateKey,
		"-----END RSA PRIVATE KEY-----",
		"</key>",
		"",
	}...)

	return lines
}

func (h *hideMyAss) PortForward(ctx context.Context, client *http.Client,
	openFile os.OpenFileFunc, pfLogger logging.Logger, gateway net.IP, fw firewall.Configurator,
	syncState func(port uint16) (pfFilepath string)) {
	panic("port forwarding is not supported for hideMyAss")
}
