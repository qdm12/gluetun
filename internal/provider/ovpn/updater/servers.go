package updater

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

var ErrResponseSuccessFalse = errors.New("response success field is false")

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error,
) {
	data, err := fetchAPI(ctx, u.client)
	if err != nil {
		return nil, fmt.Errorf("fetching API: %w", err)
	} else if !data.Success {
		return nil, fmt.Errorf("%w", ErrResponseSuccessFalse)
	}

	for dataCenterIndex, dataCenter := range data.DataCenters {
		err = dataCenter.validate()
		if err != nil {
			return nil, fmt.Errorf("validating data center %d of %d: %w",
				dataCenterIndex+1, len(data.DataCenters), err)
		}

		for _, apiServer := range dataCenter.Servers {
			if !apiServer.Online {
				continue
			}

			baseServer := models.Server{
				Country:  dataCenter.CountryName,
				City:     dataCenter.City,
				Hostname: apiServer.Ptr,
				IPs:      []netip.Addr{apiServer.IP},
			}
			openVPNServer := baseServer
			openVPNServer.VPN = vpn.OpenVPN
			openVPNServer.TCP = true
			openVPNServer.UDP = true
			servers = append(servers, openVPNServer)

			wireguardServer := baseServer
			wireguardServer.VPN = vpn.Wireguard
			wireguardServer.WgPubKey = apiServer.PublicKey
			servers = append(servers, wireguardServer)
		}
	}

	if len(servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(servers), minServers)
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
