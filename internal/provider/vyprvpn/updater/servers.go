// Package vyprvpn contains code to obtain the server information
// for the VyprVPN provider.
package vyprvpn

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	const url = "https://support.vyprvpn.com/hc/article_attachments/360052617332/Vypr_OpenVPN_20200320.zip"
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
			continue // not an OpenVPN file
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

		tcp, udp, err := openvpn.ExtractProto(content)
		if err != nil {
			// treat error as warning and go to next file
			u.warner.Warn(err.Error() + " in " + fileName)
			continue
		}

		region, err := parseFilename(fileName)
		if err != nil {
			// treat error as warning and go to next file
			u.warner.Warn(err.Error())
			continue
		}

		hts.add(host, region, tcp, udp)
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

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
