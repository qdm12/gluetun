// Package purevpn contains code to obtain the server information
// for the PureVPN provider.
package purevpn

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	const url = "https://d32d3g1fvkpl8y.cloudfront.net/heartbleed/windows/New+OVPN+Files.zip"
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
	hostToIPs, warnings, err := resolveHosts(ctx, u.presolver, hosts, minServers)
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
	if err != nil {
		return nil, err
	}

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	// Get public IP address information
	ipsToGetInfo := make([]net.IP, len(servers))
	for i := range servers {
		ipsToGetInfo[i] = servers[i].IPs[0]
	}
	ipsInfo, err := publicip.MultiInfo(ctx, u.client, ipsToGetInfo)
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
