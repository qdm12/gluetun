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

	hosts := make(map[string]struct{}, len(data.Servers))

	for _, serverData := range data.Servers {
		openVPNHost := serverData.Hostnames.OpenVPN
		if openVPNHost != "" {
			hosts[openVPNHost] = struct{}{}
		}

		wireguardHost := serverData.Hostnames.Wireguard
		if wireguardHost != "" {
			hosts[wireguardHost] = struct{}{}
		}
	}

	if len(hosts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hosts), minServers)
	}

	hostsSlice := make(sort.StringSlice, 0, len(hosts))
	for host := range hosts {
		hostsSlice = append(hostsSlice, host)
	}
	hostsSlice.Sort() // for predictable unit tests

	resolveSettings := parallelResolverSettings(hostsSlice)
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

	servers = make([]models.Server, 0, len(hostToIPs))
	for _, serverData := range data.Servers {
		server := models.Server{
			Country: serverData.Country,
			City:    serverData.City,
			ISP:     serverData.ISP,
		}

		openVPNHostname := serverData.Hostnames.OpenVPN
		wireguardHostname := serverData.Hostnames.Wireguard
		if openVPNHostname == "" && wireguardHostname == "" {
			warning := fmt.Sprintf("server data %v has no OpenVPN nor Wireguard hostname", serverData)
			warnings = append(warnings, warning)
			continue
		}

		if openVPNHostname != "" {
			openVPNServer := server
			openVPNServer.Hostname = openVPNHostname
			openVPNServer.VPN = vpn.OpenVPN
			openVPNServer.UDP = true
			openVPNServer.TCP = true
			openVPNServer.IPs = hostToIPs[openVPNHostname]
			servers = append(servers, openVPNServer)
		}

		if wireguardHostname != "" {
			wireguardServer := server
			wireguardServer.Hostname = wireguardHostname
			wireguardServer.VPN = vpn.Wireguard
			wireguardServer.IPs = hostToIPs[wireguardHostname]
			wireguardServer.WgPubKey = serverData.WgPubKey
			servers = append(servers, wireguardServer)
		}
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
