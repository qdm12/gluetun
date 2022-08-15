package updater

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	hostToData, err := fetchServers(ctx, u.client)
	if err != nil {
		return nil, fmt.Errorf("fetching and parsing website: %w", err)
	}

	openvpnURLs := make([]string, 0, len(hostToData))
	for _, data := range hostToData {
		openvpnURLs = append(openvpnURLs, data.ovpnURL)
	}

	if len(openvpnURLs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(openvpnURLs), minServers)
	}

	const failEarly = false // some URLs from the website are not valid
	udpHostToURL, errors := openvpn.FetchMultiFiles(ctx, u.client, openvpnURLs, failEarly)
	for _, err := range errors {
		u.warner.Warn(fmt.Sprintf("fetching OpenVPN files: %s", err))
	}

	if len(udpHostToURL) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(udpHostToURL), minServers)
	}

	hosts := make([]string, 0, len(udpHostToURL))
	for host := range udpHostToURL {
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
			common.ErrNotEnoughServers, len(hosts), minServers)
	}

	servers = make([]models.Server, 0, len(hostToIPs))
	for host, IPs := range hostToIPs {
		_, udp := udpHostToURL[host]

		serverData := hostToData[host]

		server := models.Server{
			VPN:      vpn.OpenVPN,
			Region:   serverData.region,
			Country:  serverData.country,
			City:     serverData.city,
			Hostname: host,
			UDP:      udp,
			IPs:      IPs,
		}
		servers = append(servers, server)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
