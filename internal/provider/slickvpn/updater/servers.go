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
	servers []models.Server, err error,
) {
	// Since SlickVPN website listing VPN servers https://www.slickvpn.com/locations/
	// went to become a pile of trash, we now hardcode the servers data below.
	servers = []models.Server{
		{Hostname: "gw2.sin2.slickvpn.com", Region: "Asia", Country: "Singapore", City: "Singapore"},
		{Hostname: "gw1.bos1.slickvpn.com", Region: "North America", Country: "United States", City: "Boston"},
		{Hostname: "gw1.cmh1.slickvpn.com", Region: "North America", Country: "United States", City: "Columbus"},
		{Hostname: "gw1.lax2.slickvpn.com", Region: "North America", Country: "United States", City: "Los Angeles"},
		{Hostname: "gw1.lga2.slickvpn.com", Region: "North America", Country: "United States", City: "New York"},
		{Hostname: "gw1.man2.slickvpn.com", Region: "Europe", Country: "United Kingdom", City: "Manchester"},
		{Hostname: "gw2.ams3.slickvpn.com", Region: "Europe", Country: "Netherlands", City: "Amsterdam"},
		{Hostname: "gw2.hou1.slickvpn.com", Region: "North America", Country: "United States", City: "Houston"},
		{Hostname: "gw1.mci2.slickvpn.com", Region: "North America", Country: "United States", City: "Kansas City"},
		{Hostname: "gw2.slc1.slickvpn.com", Region: "North America", Country: "United States", City: "Salt Lake City"},
		{Hostname: "gw1.stl1.slickvpn.com", Region: "North America", Country: "United States", City: "St Louis"},
	}

	hosts := make([]string, len(servers))
	for i := range servers {
		hosts[i] = servers[i].Hostname
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

	for i := range servers {
		servers[i].VPN = vpn.OpenVPN
		servers[i].TCP = true
		servers[i].UDP = true
		servers[i].IPs = hostToIPs[servers[i].Hostname]
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
