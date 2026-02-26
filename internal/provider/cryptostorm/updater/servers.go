package updater

import (
	"context"
	"fmt"
	"net/netip"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error,
) {
	nodes, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		if node.IPv4 == "" {
			continue
		}

		ip, err := netip.ParseAddr(node.IPv4)
		if err != nil {
			return nil, fmt.Errorf("parsing IP for node %s: %w", node.Hostname, err)
		}

		// WireGuard server entry (only if public key is available).
		if node.WgPubKey != "" {
			server := models.Server{
				VPN:         vpn.Wireguard,
				Country:     node.Country,
				City:        node.City,
				Hostname:    node.Hostname,
				WgPubKey:    node.WgPubKey,
				PortForward: node.PortFwd,
				IPs:         []netip.Addr{ip},
			}
			servers = append(servers, server)
		}

		// OpenVPN server entry.
		// Derive x509 name from hostname: "newyork.cstorm.is" -> "cryptostorm newyork server"
		location := strings.Split(node.Hostname, ".")[0]
		ovpnX509 := "cryptostorm " + location + " server"
		openvpnServer := models.Server{
			VPN:         vpn.OpenVPN,
			Country:     node.Country,
			City:        node.City,
			Hostname:    node.Hostname,
			OvpnX509:    ovpnX509,
			TCP:         true,
			UDP:         true,
			PortForward: node.PortFwd,
			IPs:         []netip.Addr{ip},
		}
		servers = append(servers, openvpnServer)
	}

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
