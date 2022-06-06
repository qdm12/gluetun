// Package windscribe contains code to obtain the server information
// for the Windscribe provider.
package windscribe

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

var (
	ErrNoWireguardKey = errors.New("no wireguard public key found")
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	data, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, err
	}

	for _, regionData := range data.Data {
		region := regionData.Region
		for _, group := range regionData.Groups {
			city := group.City
			x5090Name := group.OvpnX509
			wgPubKey := group.WgPubKey
			for _, node := range group.Nodes {
				ips := make([]net.IP, 0, 2) // nolint:gomnd
				if node.IP != nil {
					ips = append(ips, node.IP)
				}
				if node.IP2 != nil {
					ips = append(ips, node.IP2)
				}
				server := models.Server{
					VPN:      vpn.OpenVPN,
					TCP:      true,
					UDP:      true,
					Region:   region,
					City:     city,
					Hostname: node.Hostname,
					OvpnX509: x5090Name,
					IPs:      ips,
				}
				servers = append(servers, server)

				if node.IP3 == nil { // Wireguard + Stealth
					continue
				} else if wgPubKey == "" {
					return nil, fmt.Errorf("%w: for node %s", ErrNoWireguardKey, node.Hostname)
				}

				server.VPN = vpn.Wireguard
				server.UDP = true
				server.TCP = false
				server.OvpnX509 = ""
				server.WgPubKey = wgPubKey
				server.IPs = []net.IP{node.IP3}
				servers = append(servers, server)
			}
		}
	}

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
