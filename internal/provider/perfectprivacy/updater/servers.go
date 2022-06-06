package perfectprivacy

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/openvpn"
)

func (u *Updater) GetServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
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

	contents, err := u.unzipper.FetchAndExtract(ctx, zipURL.String())
	if err != nil {
		return nil, err
	}

	cts := make(cityToServer)

	for fileName, content := range contents {
		err := addServerFromOvpn(cts, fileName, content)
		if err != nil {
			u.warner.Warn(err.Error() + " in " + fileName)
		}
	}

	if len(cts) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(cts), minServers)
	}

	servers = cts.toServersSlice()

	sort.Sort(models.SortableServers(servers))

	return servers, nil
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
