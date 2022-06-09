package updater

import (
	"context"
	"net"
	"net/http"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/publicip"
)

func setLocationInfo(ctx context.Context, client *http.Client, servers []models.Server) (err error) {
	// Get public IP address information
	ipsToGetInfo := make([]net.IP, 0, len(servers))
	for _, server := range servers {
		ipsToGetInfo = append(ipsToGetInfo, server.IPs...)
	}
	ipsInfo, err := publicip.MultiInfo(ctx, client, ipsToGetInfo)
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
