package updater

import (
	"context"
	"net/netip"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/publicip/api"
)

func setLocationInfo(ctx context.Context, fetcher common.IPFetcher, servers []models.Server) (err error) {
	// Get public IP address information
	ipsToGetInfo := make([]netip.Addr, 0, len(servers))
	for _, server := range servers {
		ipsToGetInfo = append(ipsToGetInfo, server.IPs...)
	}
	ipsInfo, err := api.FetchMultiInfo(ctx, fetcher, ipsToGetInfo)
	if err != nil {
		return err
	}

	for i := range servers {
		ipInfo := ipsInfo[i]
		servers[i].Country = ipInfo.Country
		servers[i].Region = ipInfo.Region
		servers[i].City = ipInfo.City
	}

	return nil
}
