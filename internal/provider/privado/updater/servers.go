package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/netip"
	"sort"

	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
)

func (u *Updater) FetchServers(ctx context.Context, minServers int) (
	servers []models.Server, err error,
) {
	const url = "https://privadovpn.com/apps/servers_export.json"
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	response, err := u.client.Do(request)
	if err != nil {
		return nil, err
	}

	var data struct {
		Servers []struct {
			Country  string     `json:"country"`
			City     string     `json:"city"`
			Hostname string     `json:"hostname"`
			IP       netip.Addr `json:"ip"`
		} `json:"servers"`
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&data)
	if err != nil {
		_ = response.Body.Close()
		return nil, fmt.Errorf("decoding JSON response: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing response body: %w", err)
	}

	if len(data.Servers) < minServers {
		return nil, fmt.Errorf("%w: %d and expected at least %d",
			common.ErrNotEnoughServers, len(data.Servers), minServers)
	}

	servers = make([]models.Server, len(data.Servers))
	for i, server := range data.Servers {
		servers[i] = models.Server{
			VPN:      vpn.OpenVPN,
			Country:  server.Country,
			City:     server.City,
			Hostname: server.Hostname,
			IPs:      []netip.Addr{server.IP},
			UDP:      true,
		}
	}

	sort.Sort(models.SortableServers(servers))

	return servers, nil
}
