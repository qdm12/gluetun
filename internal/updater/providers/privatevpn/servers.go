// Package privatevpn contains code to obtain the server information
// for the PrivateVPN provider.
package privatevpn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, unzipper unzip.Unzipper,
	presolver resolver.Parallel, minServers int) (
	servers []models.Server, warnings []string, err error) {
	const url = "https://privatevpn.com/client/PrivateVPN-TUN.zip"
	contents, err := unzipper.FetchAndExtract(ctx, url)
	if err != nil {
		return nil, nil, err
	} else if len(contents) < minServers {
		return nil, nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(contents), minServers)
	}

	countryCodes := constants.CountryCodes()

	hts := make(hostToServer)
	noHostnameServers := make([]models.Server, 0, 1) // there is only one for now

	for fileName, content := range contents {
		if !strings.HasSuffix(fileName, ".ovpn") {
			continue // not an OpenVPN file
		}

		countryCode, city, err := parseFilename(fileName)
		if err != nil {
			warnings = append(warnings, err.Error())
			continue
		}

		country, warning := codeToCountry(countryCode, countryCodes)
		if warning != "" {
			warnings = append(warnings, warning)
		}

		host, warning, err := openvpn.ExtractHost(content)
		if warning != "" {
			warnings = append(warnings, warning)
		}
		if err == nil { // found host
			hts.add(host, country, city)
			continue
		}

		ips, extractIPErr := openvpn.ExtractIPs(content)
		if warning != "" {
			warnings = append(warnings, warning)
		}
		if extractIPErr != nil {
			// treat extract host error as warning and go to next file
			warning := err.Error() + " in " + fileName
			warnings = append(warnings, warning)
			continue
		}
		server := models.Server{
			Country: country,
			City:    city,
			IPs:     ips,
		}
		noHostnameServers = append(noHostnameServers, server)
	}

	if len(noHostnameServers)+len(hts) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(servers)+len(hts), minServers)
	}

	hosts := hts.toHostsSlice()

	hostToIPs, newWarnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	hts.adaptWithIPs(hostToIPs)

	servers = hts.toServersSlice()
	servers = append(servers, noHostnameServers...)

	sortServers(servers)

	return servers, warnings, nil
}
