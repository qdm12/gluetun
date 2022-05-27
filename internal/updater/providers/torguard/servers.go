// Package torguard contains code to obtain the server information
// for the Torguard provider.
package torguard

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/qdm12/gluetun/internal/updater/unzip"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, unzipper unzip.Unzipper,
	presolver resolver.Parallel, minServers int) (
	servers []models.Server, warnings []string, err error) {
	const tcpURL = "https://torguard.net/downloads/OpenVPN-TCP-Linux.zip"
	tcpContents, err := unzipper.FetchAndExtract(ctx, tcpURL)
	if err != nil {
		return nil, nil, err
	}

	const udpURL = "https://torguard.net/downloads/OpenVPN-UDP-Linux.zip"
	udpContents, err := unzipper.FetchAndExtract(ctx, udpURL)
	if err != nil {
		return nil, nil, err
	}

	hts := make(hostToServer)
	titleCaser := cases.Title(language.English)

	for fileName, content := range tcpContents {
		const tcp, udp = true, false
		newWarnings := addServerFromOvpn(fileName, content, hts, tcp, udp, titleCaser)
		warnings = append(warnings, newWarnings...)
	}

	for fileName, content := range udpContents {
		const tcp, udp = false, true
		newWarnings := addServerFromOvpn(fileName, content, hts, tcp, udp, titleCaser)
		warnings = append(warnings, newWarnings...)
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

func addServerFromOvpn(fileName string, content []byte,
	hts hostToServer, tcp, udp bool, titleCaser cases.Caser) (warnings []string) {
	if !strings.HasSuffix(fileName, ".ovpn") {
		return nil // not an OpenVPN file
	}

	country, city := parseFilename(fileName, titleCaser)

	host, warning, err := openvpn.ExtractHost(content)
	if warning != "" {
		warnings = append(warnings, warning)
	}
	if err != nil {
		// treat error as warning and go to next file
		warning := err.Error() + " in " + fileName
		warnings = append(warnings, warning)
		return warnings
	}

	ips, err := openvpn.ExtractIPs(content)
	if err != nil {
		// treat error as warning and go to next file
		warning := err.Error() + " in " + fileName
		warnings = append(warnings, warning)
		return warnings
	}

	hts.add(host, country, city, tcp, udp, ips)
	return warnings
}
