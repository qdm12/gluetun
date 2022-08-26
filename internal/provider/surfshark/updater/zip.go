package updater

import (
	"context"
	"strings"

	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/provider/surfshark/servers"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func addOpenVPNServersFromZip(ctx context.Context,
	unzipper common.Unzipper, hts hostToServers) (
	warnings []string, err error) {
	const url = "https://my.surfshark.com/vpn/api/v1/server/configurations"
	contents, err := unzipper.FetchAndExtract(ctx, url)
	if err != nil {
		return nil, err
	}

	hostnamesDone := hts.toHostsSlice()
	hostnamesDoneSet := make(map[string]struct{}, len(hostnamesDone))
	for _, hostname := range hostnamesDone {
		hostnamesDoneSet[hostname] = struct{}{}
	}

	locationData := servers.LocationData()
	hostToLocation := hostToLocation(locationData)

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
			// TODO gather location data for IP address Openvpn files
			// and process those when this error triggers.
			warnings = append(warnings, warning)
			continue
		}

		_, ok := hostnamesDoneSet[host]
		if ok {
			continue // already done in API
		}

		tcp, udp, err := openvpn.ExtractProto(content)
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error() + " in " + fileName
			warnings = append(warnings, warning)
			continue
		}

		data, err := getHostInformation(host, hostToLocation)
		if err != nil {
			// treat error as warning and go to next file
			warning := err.Error()
			warnings = append(warnings, warning)
			continue
		}

		hts.addOpenVPN(host, data.Region, data.Country, data.City,
			data.RetroLoc, tcp, udp)
	}

	return warnings, nil
}
