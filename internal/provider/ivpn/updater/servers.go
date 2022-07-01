package updater

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	data, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, fmt.Errorf("failed fetching API: %w", err)
	}

	hosts := make([]string, 0, len(data.Servers))

	for _, serverData := range data.Servers {
		openVPNHost := serverData.Hostnames.OpenVPN
		if openVPNHost != "" {
			hosts = append(hosts, openVPNHost)
		}

		wireguardHost := serverData.Hostnames.Wireguard
		if wireguardHost != "" {
			hosts = append(hosts, wireguardHost)
		}
	}

	if len(hosts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hosts), minServers)
	}

	resolveSettings := parallelResolverSettings(hosts)
	hostToIPs, warnings, err := u.parallelResolver.Resolve(ctx, resolveSettings)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, err
	}

	if len(hostToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	servers = make([]models.Server, 0, len(hosts))
	for _, serverData := range data.Servers {
		vpnType := vpn.OpenVPN
		hostname := serverData.Hostnames.OpenVPN
		tcp := true
		wgPubKey := ""
		if hostname == "" {
			vpnType = vpn.Wireguard
			hostname = serverData.Hostnames.Wireguard
			tcp = false
			wgPubKey = serverData.WgPubKey
		}

		server := models.Server{
			VPN:      vpnType,
			Country:  serverData.Country,
			City:     serverData.City,
			ISP:      serverData.ISP,
			Hostname: hostname,
			WgPubKey: wgPubKey,
			TCP:      tcp,
			UDP:      true,
			IPs:      hostToIPs[hostname],
		}
		servers = append(servers, server)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
