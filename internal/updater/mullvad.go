package updater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateMullvad() (err error) {
	servers, err := findMullvadServers(u.httpGet)
	if err != nil {
		return fmt.Errorf("cannot update Mullvad servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyMullvadServers(servers))
	}
	u.servers.Mullvad.Timestamp = u.timeNow().Unix()
	u.servers.Mullvad.Servers = servers
	return nil
}

func findMullvadServers(httpGet httpGetFunc) (servers []models.MullvadServer, err error) {
	const url = "https://api.mullvad.net/www/relays/openvpn/"
	response, err := httpGet(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var data []struct {
		Country  string `json:"country_name"`
		City     string `json:"city_name"`
		Active   bool   `json:"active"`
		Owned    bool   `json:"owned"`
		Provider string `json:"provider"`
		IPv4     string `json:"ipv4_addr_in"`
		IPv6     string `json:"ipv6_addr_in"`
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}
	serversByKey := map[string]models.MullvadServer{}
	for _, jsonServer := range data {
		if !jsonServer.Active {
			continue
		}
		ipv4 := net.ParseIP(jsonServer.IPv4)
		ipv6 := net.ParseIP(jsonServer.IPv6)
		if ipv4 == nil || ipv4.To4() == nil {
			return nil, fmt.Errorf("cannot parse ipv4 address %q", jsonServer.IPv4)
		} else if ipv6 == nil || ipv6.To4() != nil {
			return nil, fmt.Errorf("cannot parse ipv6 address %q", jsonServer.IPv6)
		}
		key := fmt.Sprintf("%s%s%t%s", jsonServer.Country, jsonServer.City, jsonServer.Owned, jsonServer.Provider)
		if server, ok := serversByKey[key]; ok {
			server.IPs = append(server.IPs, ipv4)
			server.IPsV6 = append(server.IPsV6, ipv6)
			serversByKey[key] = server
		} else {
			serversByKey[key] = models.MullvadServer{
				IPs:     []net.IP{ipv4},
				IPsV6:   []net.IP{ipv6},
				Country: jsonServer.Country,
				City:    jsonServer.City,
				ISP:     jsonServer.Provider,
				Owned:   jsonServer.Owned,
			}
		}
	}
	for _, server := range serversByKey {
		server.IPs = uniqueSortedIPs(server.IPs)
		server.IPsV6 = uniqueSortedIPs(server.IPsV6)
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Country+servers[i].City+servers[i].ISP < servers[j].Country+servers[j].City+servers[j].ISP
	})
	return servers, nil
}

//nolint:goconst
func stringifyMullvadServers(servers []models.MullvadServer) (s string) {
	s = "func MullvadServers() []models.MullvadServer {\n"
	s += "	return []models.MullvadServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
