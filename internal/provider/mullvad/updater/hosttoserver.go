package updater

import (
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

type hostToServer map[string]models.Server

var (
	ErrNoIP                = errors.New("no IP address for VPN server")
	ErrParseIPv4           = errors.New("cannot parse IPv4 address")
	ErrParseIPv6           = errors.New("cannot parse IPv6 address")
	ErrVPNTypeNotSupported = errors.New("VPN type not supported")
)

func (hts hostToServer) add(data serverData) (err error) {
	if !data.Active {
		return
	}

	if data.IPv4 == "" && data.IPv6 == "" {
		return ErrNoIP
	}

	server, ok := hts[data.Hostname]
	if ok { // API returns a server per hostname at most
		return nil
	}

	switch data.Type {
	case "openvpn":
		server.VPN = vpn.OpenVPN
		server.UDP = true
		server.TCP = true
	case "wireguard":
		server.VPN = vpn.Wireguard
	case "bridge":
		// ignore bridge servers
		return nil
	default:
		return fmt.Errorf("%w: %s", ErrVPNTypeNotSupported, data.Type)
	}

	if data.IPv4 != "" {
		ipv4 := net.ParseIP(data.IPv4)
		if ipv4 == nil || ipv4.To4() == nil {
			return fmt.Errorf("%w: %s", ErrParseIPv4, data.IPv4)
		}
		server.IPs = append(server.IPs, ipv4)
	}

	if data.IPv6 != "" {
		ipv6 := net.ParseIP(data.IPv6)
		if ipv6 == nil || ipv6.To4() != nil {
			return fmt.Errorf("%w: %s", ErrParseIPv6, data.IPv6)
		}
		server.IPs = append(server.IPs, ipv6)
	}

	server.Country = data.Country
	server.City = strings.ReplaceAll(data.City, ",", "")
	server.Hostname = data.Hostname
	server.ISP = data.Provider
	server.Owned = data.Owned
	server.WgPubKey = data.PubKey

	hts[data.Hostname] = server

	return nil
}

func (hts hostToServer) toServersSlice() (servers []models.Server) {
	servers = make([]models.Server, 0, len(hts))
	for _, server := range hts {
		server.IPs = uniqueSortedIPs(server.IPs)
		servers = append(servers, server)
	}
	return servers
}
