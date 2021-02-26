package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateWindscribe(ctx context.Context) (err error) {
	servers, err := findWindscribeServers(ctx, u.client)
	if err != nil {
		return fmt.Errorf("cannot update Windscribe servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyWindscribeServers(servers))
	}
	u.servers.Windscribe.Timestamp = u.timeNow().Unix()
	u.servers.Windscribe.Servers = servers
	return nil
}

func findWindscribeServers(ctx context.Context, client *http.Client) (servers []models.WindscribeServer, err error) {
	const baseURL = "https://assets.windscribe.com/serverlist/mob-v2/1/"
	cacheBreaker := time.Now().Unix()
	url := fmt.Sprintf("%s%d", baseURL, cacheBreaker)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %s", ErrHTTPStatusCodeNotOK, response.Status)
	}

	decoder := json.NewDecoder(response.Body)
	var jsonData struct {
		Data []struct {
			Region string `json:"name"`
			Groups []struct {
				City  string `json:"city"`
				Nodes []struct {
					Hostname  string `json:"hostname"`
					OpenvpnIP net.IP `json:"ip2"`
				} `json:"nodes"`
			} `json:"groups"`
		} `json:"data"`
	}
	if err := decoder.Decode(&jsonData); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUnmarshalResponseBody, err)
	}

	if err := response.Body.Close(); err != nil {
		return nil, err
	}

	for _, regionBlock := range jsonData.Data {
		region := regionBlock.Region
		for _, group := range regionBlock.Groups {
			city := group.City
			for _, node := range group.Nodes {
				server := models.WindscribeServer{
					Region:   region,
					City:     city,
					Hostname: node.Hostname,
					IP:       node.OpenvpnIP,
				}
				servers = append(servers, server)
			}
		}
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region+servers[i].City+servers[i].Hostname <
			servers[j].Region+servers[j].City+servers[j].Hostname
	})
	return servers, nil
}

func stringifyWindscribeServers(servers []models.WindscribeServer) (s string) {
	s = "func WindscribeServers() []models.WindscribeServer {\n"
	s += "	return []models.WindscribeServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
