// Package surfshark contains code to obtain the server information
// for the Surshark provider.
package surfshark

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, unzipper unzip.Unzipper,
	presolver resolver.Parallel, minServers int) (
	servers []models.SurfsharkServer, warnings []string, err error) {
	const url = "https://my.surfshark.com/vpn/api/v1/server/configurations"
	contents, err := unzipper.FetchAndExtract(ctx, url)
	if err != nil {
		return nil, nil, err
	} else if len(contents) < minServers {
		return nil, nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(contents), minServers)
	}

	subdomainToRegion := subdomainToRegion()
	hts := make(hostToServer)

	for fileName, content := range contents {
		if !strings.HasSuffix(fileName, ".ovpn") {
			continue // not an OpenVPN file
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

		region, err := parseHost(host, subdomainToRegion)
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error()
			warnings = append(warnings, warning)
			continue
		}

		hts.add(host, region)
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

	// process subdomain entries in mapping that were not in the Zip file
	subdomainsDone := hts.toSubdomainsSlice()
	for _, subdomainDone := range subdomainsDone {
		delete(subdomainToRegion, subdomainDone)
	}
	remainingServers, newWarnings, err := getRemainingServers(
		ctx, subdomainToRegion, presolver)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	servers = append(servers, remainingServers...)

	if len(servers) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers), minServers)
	}

	sortServers(servers)

	return servers, warnings, nil
}

func getRemainingServers(ctx context.Context,
	subdomainToRegionLeft map[string]string, presolver resolver.Parallel) (
	servers []models.SurfsharkServer, warnings []string, err error) {
	hosts := make([]string, 0, len(subdomainToRegionLeft))
	const suffix = ".prod.surfshark.com"
	for subdomain := range subdomainToRegionLeft {
		hosts = append(hosts, subdomain+suffix)
	}

	const minServers = 0
	hostToIPs, warnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	if err != nil {
		return nil, warnings, err
	}

	servers = make([]models.SurfsharkServer, 0, len(hostToIPs))
	for host, IPs := range hostToIPs {
		region, err := parseHost(host, subdomainToRegionLeft)
		if err != nil {
			return nil, warnings, err
		}
		server := models.SurfsharkServer{
			Region: region,
			IPs:    IPs,
		}
		servers = append(servers, server)
	}

	return servers, warnings, nil
}
