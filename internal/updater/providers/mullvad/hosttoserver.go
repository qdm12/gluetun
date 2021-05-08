package mullvad

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.MullvadServer

var (
	ErrParseIPv4 = errors.New("cannot parse IPv4 address")
	ErrParseIPv6 = errors.New("cannot parse IPv6 address")
)

func (hts hostToServer) add(data serverData) (err error) {
	if !data.Active {
		return
	}

	ipv4 := net.ParseIP(data.IPv4)
	if ipv4 == nil || ipv4.To4() == nil {
		return fmt.Errorf("%w: %s", ErrParseIPv4, data.IPv4)
	}

	ipv6 := net.ParseIP(data.IPv6)
	if ipv6 == nil || ipv6.To4() != nil {
		return fmt.Errorf("%w: %s", ErrParseIPv6, data.IPv6)
	}

	server, ok := hts[data.Hostname]
	if !ok {
		server.Country = data.Country
		server.City = strings.ReplaceAll(data.City, ",", "")
		server.ISP = data.Provider
		server.Owned = data.Owned
	}

	server.IPs = append(server.IPs, ipv4)
	server.IPsV6 = append(server.IPsV6, ipv6)

	hts[data.Hostname] = server

	return nil
}

func (hts hostToServer) toServersSlice() (servers []models.MullvadServer) {
	servers = make([]models.MullvadServer, 0, len(hts))
	for _, server := range hts {
		server.IPs = uniqueSortedIPs(server.IPs)
		server.IPsV6 = uniqueSortedIPs(server.IPsV6)
		servers = append(servers, server)
	}
	return servers
}
