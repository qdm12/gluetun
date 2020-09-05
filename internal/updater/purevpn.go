package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePurevpn(ctx context.Context) (err error) {
	servers, warnings, err := findPurevpnServers(ctx, u.httpGet, u.lookupIP)
	for _, warning := range warnings {
		u.println(warning)
	}
	if err != nil {
		return fmt.Errorf("cannot update Purevpn servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyPurevpnServers(servers))
	}
	u.servers.Purevpn.Timestamp = u.timeNow().Unix()
	u.servers.Purevpn.Servers = servers
	return nil
}

func findPurevpnServers(ctx context.Context, httpGet httpGetFunc, lookupIP lookupIPFunc) (
	servers []models.PurevpnServer, warnings []string, err error) {
	const url = "https://support.purevpn.com/vpn-servers"
	response, err := httpGet(url)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf(response.Status)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, nil, err
	}
	const jsonPrefix = "<script>var servers = "
	const jsonSuffix = "</script>"
	s := string(bytes)
	jsonPrefixIndex := strings.Index(s, jsonPrefix)
	if jsonPrefixIndex == -1 {
		return nil, nil, fmt.Errorf("cannot find %q in html", jsonPrefix)
	}
	s = s[jsonPrefixIndex+len(jsonPrefix):]
	endIndex := strings.Index(s, jsonSuffix)
	if endIndex == -1 {
		return nil, nil, fmt.Errorf("cannot find %q after %q in html", jsonSuffix, jsonPrefix)
	}
	s = s[:endIndex]
	var data []struct {
		Region  string `json:"region_name"`
		Country string `json:"country_name"`
		City    string `json:"city_name"`
		TCP     string `json:"tcp"`
		UDP     string `json:"udp"`
	}
	if err := json.Unmarshal([]byte(s), &data); err != nil {
		return nil, nil, err
	}
	sort.Slice(data, func(i, j int) bool {
		if data[i].Region == data[j].Region {
			if data[i].Country == data[j].Country {
				return data[i].City < data[j].City
			}
			return data[i].Country < data[j].Country
		}
		return data[i].Region < data[j].Region
	})
	for _, jsonServer := range data {
		if jsonServer.UDP == "" && jsonServer.TCP == "" {
			warnings = append(warnings, fmt.Sprintf("server %s %s %s does not support TCP and UDP for openvpn", jsonServer.Region, jsonServer.Country, jsonServer.City))
			continue
		}
		if jsonServer.UDP == "" || jsonServer.TCP == "" {
			warnings = append(warnings, fmt.Sprintf("server %s %s %s does not support TCP or UDP for openvpn", jsonServer.Region, jsonServer.Country, jsonServer.City))
			continue
		}
		host := jsonServer.UDP
		const repetition = 3
		IPs, err := resolveRepeat(ctx, lookupIP, host, repetition)
		if err != nil {
			warnings = append(warnings, err.Error())
			continue
		}
		servers = append(servers, models.PurevpnServer{
			Region:  jsonServer.Region,
			Country: jsonServer.Country,
			City:    jsonServer.City,
			IPs:     IPs,
		})
	}
	return servers, warnings, nil
}

func stringifyPurevpnServers(servers []models.PurevpnServer) (s string) {
	s = "func PurevpnServers() []models.PurevpnServer {\n"
	s += "	return []models.PurevpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
