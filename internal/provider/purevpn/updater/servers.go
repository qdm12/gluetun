package updater

import (
	"context"
	"fmt"
	"net/netip"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	const url = "https://d11a57lttb2ffq.cloudfront.net/heartbleed/router/Recommended-CA2.zip"
	contents, err := u.unzipper.FetchAndExtract(ctx, url)
	if err != nil {
		return nil, err
	} else if len(contents) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(contents), minServers)
	}

	hts := make(hostToServer)

	for fileName, content := range contents {
		if !strings.HasSuffix(fileName, ".ovpn") {
			continue
		}

		tcp, udp, err := openvpn.ExtractProto(content)
		if err != nil {
			// treat error as warning and go to next file
			u.warner.Warn(err.Error() + " in " + fileName)
			continue
		}

		host, warning, err := openvpn.ExtractHost(content)
		if warning != "" {
			u.warner.Warn(warning)
		}

		if err != nil {
			// treat error as warning and go to next file
			u.warner.Warn(err.Error() + " in " + fileName)
			continue
		}

		hts.add(host, tcp, udp)
	}

	if len(hts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hts), minServers)
	}

	hosts := hts.toHostsSlice()
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

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()

	// Get public IP address information
	ipsToGetInfo := make([]netip.Addr, len(servers))
	for i := range servers {
		ipsToGetInfo[i] = servers[i].IPs[0]
	}
	ipsInfo, err := u.ipFetcher.FetchMultiInfo(ctx, ipsToGetInfo)
	if err != nil {
		return nil, err
	}
	for i := range servers {
		servers[i].Country = ipsInfo[i].Country
		servers[i].Region = ipsInfo[i].Region
		servers[i].City = ipsInfo[i].City
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
