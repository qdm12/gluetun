// Package ipvanish contains code to obtain the server information
// for the Surshark provider.
package ipvanish

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
	servers []models.Server, warnings []string, err error) {
	const url = "https://www.ipvanish.com/software/configs/configs.zip"
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
			continue // not an OpenVPN file
		}

		tcp, udp, err := openvpn.ExtractProto(content)
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error() + ": in " + fileName
			warnings = append(warnings, warning)
			continue
		}

		hostname, warning, err := openvpn.ExtractHost(content)
		if warning != "" {
			warnings = append(warnings, warning)
		}
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error() + " in " + fileName
			warnings = append(warnings, warning)
			continue
		}

		country, city, err := parseFilename(fileName, hostname)
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error() + " in " + fileName
			warnings = append(warnings, warning)
			continue
		}

		hts.add(hostname, country, city, tcp, udp)
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

	sortServers(servers)

	return servers, warnings, nil
}
