package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateNordvpn(ctx context.Context) (err error) {
	servers, warnings, err := findNordvpnServers(ctx, u.client)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Nordvpn: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Nordvpn servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyNordvpnServers(servers))
	}
	u.servers.Nordvpn.Timestamp = u.timeNow().Unix()
	u.servers.Nordvpn.Servers = servers
	return nil
}

var (
	ErrNoIDInServerName      = errors.New("no ID in server name")
	ErrInvalidIDInServerName = errors.New("invalid ID in server name")
)

func findNordvpnServers(ctx context.Context, client *http.Client) (
	servers []models.NordvpnServer, warnings []string, err error) {
	const url = "https://nordvpn.com/api/server"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("%w: %s for %s", ErrHTTPStatusCodeNotOK, response.Status, url)
	}

	decoder := json.NewDecoder(response.Body)
	var data []struct {
		IPAddress string `json:"ip_address"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Features  struct {
			UDP bool `json:"openvpn_udp"`
			TCP bool `json:"openvpn_tcp"`
		} `json:"features"`
	}
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, err
	}

	if err := response.Body.Close(); err != nil {
		return nil, nil, err
	}

	sort.Slice(data, func(i, j int) bool {
		if data[i].Country == data[j].Country {
			return data[i].Name < data[j].Name
		}
		return data[i].Country < data[j].Country
	})

	for _, jsonServer := range data {
		if !jsonServer.Features.TCP && !jsonServer.Features.UDP {
			warnings = append(warnings, fmt.Sprintf("server %q does not support TCP and UDP for openvpn", jsonServer.Name))
			continue
		}
		ip := net.ParseIP(jsonServer.IPAddress)
		if ip == nil || ip.To4() == nil {
			return nil, nil,
				fmt.Errorf("IP address %q is not a valid IPv4 address for server %q",
					jsonServer.IPAddress, jsonServer.Name)
		}
		i := strings.IndexRune(jsonServer.Name, '#')
		if i < 0 {
			return nil, nil, fmt.Errorf("%w: %s", ErrNoIDInServerName, jsonServer.Name)
		}
		idString := jsonServer.Name[i+1:]
		idUint64, err := strconv.ParseUint(idString, 10, 16)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s", ErrInvalidIDInServerName, jsonServer.Name)
		}
		server := models.NordvpnServer{
			Region: jsonServer.Country,
			Number: uint16(idUint64),
			IP:     ip,
			TCP:    jsonServer.Features.TCP,
			UDP:    jsonServer.Features.UDP,
		}
		servers = append(servers, server)
	}
	return servers, warnings, nil
}

func stringifyNordvpnServers(servers []models.NordvpnServer) (s string) {
	s = "func NordvpnServers() []models.NordvpnServer {\n"
	s += "	return []models.NordvpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
