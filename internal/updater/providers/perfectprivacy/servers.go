package perfectprivacy

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
	"github.com/qdm12/gluetun/internal/updater/unzip"
)

var ErrNotEnoughServers = errors.New("not enough servers found")

func GetServers(ctx context.Context, unzipper unzip.Unzipper, minServers int) (
	servers []models.Server, warnings []string, err error) {
	zipURL := url.URL{
		Scheme: "https",
		Host:   "www.perfect-privacy.com",
		Path:   "/downloads/openvpn/get",
	}
	values := make(url.Values)
	values.Set("system", "linux")
	values.Set("scope", "server")
	values.Set("filetype", "zip")
	values.Set("protocol", "udp") // all support both TCP and UDP
	zipURL.RawQuery = values.Encode()

	contents, err := unzipper.FetchAndExtract(ctx, zipURL.String())
	if err != nil {
		return nil, nil, err
	}

	cts := make(cityToServer)

	for fileName, content := range contents {
		err := addServerFromOvpn(cts, fileName, content)
		if err != nil {
			warnings = append(warnings, fileName+": "+err.Error())
		}
	}

	if len(cts) < minServers {
		return nil, warnings, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(cts), minServers)
	}

	servers = cts.toServersSlice()

	sortServers(servers)

	return servers, warnings, nil
}

func addServerFromOvpn(cts cityToServer,
	fileName string, content []byte) (err error) {
	if !strings.HasSuffix(fileName, ".conf") {
		return nil // not an OpenVPN file
	}

	ips, err := openvpn.ExtractIPs(content)
	if err != nil {
		return err
	}

	city := parseFilename(fileName)

	cts.add(city, ips)

	return nil
}
