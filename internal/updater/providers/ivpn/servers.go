// Package ivpn contains code to obtain the server information
// for the Surshark provider.
package ivpn

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

var (
	ErrFetchAPI         = errors.New("failed fetching API")
	ErrNotEnoughServers = errors.New("not enough servers found")
)

func GetServers(ctx context.Context, client *http.Client,
	presolver resolver.Parallel, minServers int) (
	servers []models.IvpnServer, warnings []string, err error) {
	data, err := fetchAPI(ctx, client)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", ErrFetchAPI, err)
	}

	hosts := make([]string, 0, len(data.Servers))

	for _, serverData := range data.Servers {
		openVPNHost := serverData.Hostnames.OpenVPN
		if openVPNHost != "" {
			hosts = append(hosts, openVPNHost)
		}

		wireguardHost := serverData.Hostnames.Wireguard
		if wireguardHost != "" {
			hosts = append(hosts, wireguardHost)
		}
	}

	if len(hosts) < minServers {
		return nil, nil, fmt.Errorf("%w: %d and expected at least %d",
			ErrNotEnoughServers, len(hosts), minServers)
	}

	hostToIPs, warnings, err := resolveHosts(ctx, presolver, hosts, minServers)
	if err != nil {
		return nil, warnings, err
	}

	servers = make([]models.IvpnServer, 0, len(hosts))
	for _, serverData := range data.Servers {
		vpnType := constants.OpenVPN
		hostname := serverData.Hostnames.OpenVPN
		tcp := true
		wgPubKey := ""
		if hostname == "" {
			vpnType = constants.Wireguard
			hostname = serverData.Hostnames.Wireguard
			tcp = false
			wgPubKey = serverData.WgPubKey
		}

		server := models.IvpnServer{
			VPN:      vpnType,
			Country:  serverData.Country,
			City:     serverData.City,
			ISP:      serverData.ISP,
			Hostname: hostname,
			WgPubKey: wgPubKey,
			TCP:      tcp,
			UDP:      true,
			IPs:      hostToIPs[hostname],
		}
		servers = append(servers, server)
	}

	sortServers(servers)

	return servers, warnings, nil
}
