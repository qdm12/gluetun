// Package privado contains code to obtain the server information
// for the Privado provider.
package privado

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, unzipper unzip.Unzipper,
	client *http.Client, presolver resolver.Parallel, minServers int) (
	servers []models.PrivadoServer, warnings []string, err error) {
	const url = "https://privado.io/apps/ovpn_configs.zip"
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

		hts.add(host)
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

	if err := setLocationInfo(ctx, client, servers); err != nil {
		return nil, warnings, err
	}

	sortServers(servers)

	return servers, warnings, nil
}
