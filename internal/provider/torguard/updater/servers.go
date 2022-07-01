package updater

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	const tcpURL = "https://torguard.net/downloads/OpenVPN-TCP-Linux.zip"
	tcpContents, err := u.unzipper.FetchAndExtract(ctx, tcpURL)
	if err != nil {
		return nil, err
	}

	const udpURL = "https://torguard.net/downloads/OpenVPN-UDP-Linux.zip"
	udpContents, err := u.unzipper.FetchAndExtract(ctx, udpURL)
	if err != nil {
		return nil, err
	}

	hts := make(hostToServer)
	titleCaser := cases.Title(language.English)

	for fileName, content := range tcpContents {
		const tcp, udp = true, false
		warnings := addServerFromOvpn(fileName, content, hts, tcp, udp, titleCaser)
		u.warnWarnings(warnings)
	}

	for fileName, content := range udpContents {
		const tcp, udp = false, true
		warnings := addServerFromOvpn(fileName, content, hts, tcp, udp, titleCaser)
		u.warnWarnings(warnings)
	}

	if len(hts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(hts), minServers)
	}

	hosts := hts.toHostsSlice()
	resolveSettings := parallelResolverSettings(hosts)
	hostToIPs, warnings, err := u.parallelResolver.Resolve(ctx, resolveSettings)
	u.warnWarnings(warnings)
	if err != nil {
		return nil, err
	}

	if len(hostToIPs) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
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

func (u *Updater) warnWarnings(warnings []string) {
	for _, warning := range warnings {
		u.warner.Warn(warning)
	}
}
