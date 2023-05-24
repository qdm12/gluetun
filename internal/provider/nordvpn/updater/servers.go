package updater

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

var (
	ErrNotIPv4 = errors.New("IP address is not IPv4")
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error) {
	const recommended = true
	const limit = 0
	data, err := fetchAPI(ctx, u.client, recommended, limit)
	if err != nil {
		return nil, err
	}

	servers = make([]models.Server, 0, len(data))

	for _, jsonServer := range data {
		if jsonServer.Status != "online" {
			u.warner.Warn(fmt.Sprintf("ignoring offline server %s", jsonServer.Name))
			continue
		}

		server := models.Server{
			Country:  jsonServer.country(),
			Region:   jsonServer.region(),
			City:     jsonServer.city(),
			Hostname: jsonServer.Hostname,
			IPs:      jsonServer.ips(),
		}

		number, err := parseServerName(jsonServer.Name)
		switch {
		case errors.Is(err, ErrNoIDInServerName):
			u.warner.Warn(fmt.Sprintf("%s - leaving server number as 0", err))
		case err != nil:
			u.warner.Warn(fmt.Sprintf("failed parsing server name: %s", err))
			continue
		default: // no error
			server.Number = number
		}

		var openvpnFound bool
		openVPNServer := server // accumulate UDP+TCP technologies
		openVPNServer.VPN = vpn.OpenVPN

		for _, technology := range jsonServer.Technologies {
			switch technology.Identifier {
			case "openvpn_udp":
				openvpnFound = true
				openVPNServer.UDP = true
			case "openvpn_tcp":
				openvpnFound = true
				openVPNServer.TCP = true
			default: // Ignore other technologies
				continue
			}
		}

		if openvpnFound {
			servers = append(servers, openVPNServer)
		}
	}

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
