// Package purevpn contains code to obtain the server information
// for the PureVPN provider.
package purevpn

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, client *http.Client,
	unzipper unzip.Unzipper, presolver resolver.Parallel, minServers int) (
	servers []models.PurevpnServer, warnings []string, err error) {
	const url = "https://s3-us-west-1.amazonaws.com/heartbleed/windows/New+OVPN+Files.zip"
	contents, err := unzipper.FetchAndExtract(ctx, url)
	if err != nil {
		return nil, nil, err
	} else if len(contents) < minServers {
		return nil, nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(contents), minServers)
	}

	hts := make(hostToServer)

	for fileName, content := range contents {
		if !strings.HasSuffix(fileName, ".ovpn") {
			continue
		}

		tcp, udp, err := openvpn.ExtractProto(content)
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error() + " in " + fileName
			warnings = append(warnings, warning)
			continue
		}

		host, warning, err := openvpn.ExtractHost(content)
		if warning != "" {
			warnings = append(warnings, warning)
		}

		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error() + " in " + fileName
			warnings = append(warnings, warning)
			continue
		}

		hts.add(host, tcp, udp)
	}

	if len(hts) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(hts), minServers)
	}

	hosts := hts.toHostsSlice()
	hostToIPs, newWarnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()

	if len(servers) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	// Get public IP address information
	ipsToGetInfo := make([]net.IP, len(servers))
	for i := range servers {
		ipsToGetInfo[i] = servers[i].IPs[0]
	}
	ipsInfo, err := publicip.MultiInfo(ctx, client, ipsToGetInfo)
	if err != nil {
		return nil, warnings, err
	}
	for i := range servers {
		servers[i].Country = ipsInfo[i].Country
		servers[i].Region = ipsInfo[i].Region
		servers[i].City = ipsInfo[i].City
	}

	sortServers(servers)

	return servers, warnings, nil
}
