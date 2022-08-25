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
	// FetchServers obtains information for each VPN server
	// for the VPN service provider.
	//
	// You should aim at obtaining as much information as possible
	// for each server, such as their location information.
	// Required fields for each server are:
	// - the `VPN` protocol string field
	// - the `Hostname` string field
	// - the `IPs` IP slice field
	// - have one network protocol set, either `TCP` or `UDP`
	// - If `VPN` is `wireguard`, the `WgPubKey` field to be set
	//
	// The information obtention can be done in different ways
	// or by combining ways, depending on how the provider exposes
	// this information. Some common ones are listed below:
	//
	// - you can use u.client to fetch structured (usually JSON)
	// data of the servers from an HTTP API endpoint of the provider.
	// Example in: `internal/provider/mullvad/updater`
	// - you can use u.unzipper to download, unzip and parse a zip
	// file of OpenVPN configuration files.
	// Example in: `internal/provider/fastestvpn/updater`
	// - you can use u.parallelResolver to resolve all hostnames
	// found in parallel to obtain their corresponding IP addresses.
	// Example in: `internal/provider/fastestvpn/updater`
	//
	// The following is an example code which fetches server
	// information from an HTTP API endpoint of the provider,
	// and then resolves in parallel all hostnames to get their
	// IP addresses. You should pay attention to the following:
	// - we check multiple times we have enough servers
	// before continuing processing.
	// - hosts are deduplicated to reduce parallel resolution
	// load.
	// - servers are sorted at the end.
	//
	// Once you are done, please check all the TODO comments
	// in this package and address them.
	data, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, fmt.Errorf("fetching API: %w", err)
	}

	uniqueHosts := make(map[string]struct{}, len(data.Servers))

	for _, serverData := range data.Servers {
		if serverData.OpenVPNHostname != "" {
			uniqueHosts[serverData.OpenVPNHostname] = struct{}{}
		}

		if serverData.WireguardHostname != "" {
			uniqueHosts[serverData.WireguardHostname] = struct{}{}
		}
	}

	if len(uniqueHosts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(uniqueHosts), minServers)
	}

	hosts := make([]string, 0, len(uniqueHosts))
	for host := range uniqueHosts {
		hosts = append(hosts, host)
	}

	resolveSettings := parallelResolverSettings(hosts)
	hostToIPs, warnings, err := u.parallelResolver.Resolve(ctx, resolveSettings)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, fmt.Errorf("resolving hosts: %w", err)
	}

	if len(hostToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	maxServers := 2 * len(data.Servers) //nolint:gomnd
	servers = make([]models.Server, 0, maxServers)
	for _, serverData := range data.Servers {
		server := models.Server{
			Country:  serverData.Country,
			Region:   serverData.Region,
			City:     serverData.City,
			WgPubKey: serverData.WgPubKey,
		}
		if serverData.OpenVPNHostname != "" {
			server.VPN = vpn.OpenVPN
			server.UDP = true
			server.TCP = true
			server.Hostname = serverData.OpenVPNHostname
			server.IPs = hostToIPs[serverData.OpenVPNHostname]
			servers = append(servers, server)
		}
		if serverData.WireguardHostname != "" {
			server.VPN = vpn.Wireguard
			server.Hostname = serverData.WireguardHostname
			server.IPs = hostToIPs[serverData.WireguardHostname]
			servers = append(servers, server)
		}
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
