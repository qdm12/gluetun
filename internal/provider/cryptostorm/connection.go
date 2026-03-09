package cryptostorm

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

// GetConnection selects a server matching the given criteria and resolves
// its hostname via DNS at connection time, rather than relying on
// pre-resolved IP addresses in the server list.
func (p *Provider) GetConnection(selection settings.ServerSelection, ipv6Supported bool) (
	connection models.Connection, err error,
) {
	servers, err := p.storage.FilterServers(p.Name(), selection)
	if err != nil {
		return connection, fmt.Errorf("filtering servers: %w", err)
	}

	// Pick a random server from the filtered list.
	server := servers[rand.New(p.randSource).Intn(len(servers))] //nolint:gosec

	// Resolve the hostname at connection time.
	const resolveTimeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
	defer cancel()

	network := "ip4"
	if ipv6Supported {
		network = "ip"
	}
	ips, err := net.DefaultResolver.LookupNetIP(ctx, network, server.Hostname)
	if err != nil {
		return connection, fmt.Errorf("resolving %s: %w", server.Hostname, err)
	}
	if len(ips) == 0 {
		return connection, fmt.Errorf("no IP addresses found for %s", server.Hostname)
	}
	ip := ips[rand.New(p.randSource).Intn(len(ips))] //nolint:gosec

	// Determine protocol.
	protocol := constants.UDP
	if selection.VPN == vpn.OpenVPN && selection.OpenVPN.Protocol == constants.TCP {
		protocol = constants.TCP
	}

	// Determine port (cryptostorm accepts any port 1-65535, default 443).
	const defaultPort uint16 = 443 //nolint:mnd
	port := defaultPort
	switch selection.VPN {
	case vpn.Wireguard:
		if custom := *selection.Wireguard.EndpointPort; custom > 0 {
			port = custom
		}
	default: // OpenVPN
		if custom := *selection.OpenVPN.CustomPort; custom > 0 {
			port = custom
		}
	}

	// Allow explicit endpoint IP override.
	switch selection.VPN {
	case vpn.OpenVPN:
		if t := selection.OpenVPN.EndpointIP; t.IsValid() && !t.IsUnspecified() {
			ip = t
		}
	case vpn.Wireguard:
		if t := selection.Wireguard.EndpointIP; t.IsValid() && !t.IsUnspecified() {
			ip = t
		}
	}

	return models.Connection{
		Type:     selection.VPN,
		IP:       ip,
		Port:     port,
		Protocol: protocol,
		Hostname: server.Hostname,
		PubKey:   server.WgPubKey,
	}, nil
}
