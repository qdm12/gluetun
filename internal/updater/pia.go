package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePIA(ctx context.Context) (err error) {
	const url = "https://serverlist.piaservers.net/vpninfo/servers/v5"
	b, status, err := u.client.Get(ctx, url)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("HTTP status code %d: %s", status, strings.ReplaceAll(string(b), "\n", ""))
	}

	// remove key/signature at the bottom
	i := bytes.IndexRune(b, '\n')
	b = b[:i]

	var data struct {
		Regions []struct {
			Name        string `json:"name"`
			PortForward bool   `json:"port_forward"`
			Servers     struct {
				UDP []struct {
					IP net.IP `json:"ip"`
					CN string `json:"cn"`
				} `json:"ovpnudp"`
				TCP []struct {
					IP net.IP `json:"ip"`
					CN string `json:"cn"`
				} `json:"ovpntcp"`
			} `json:"servers"`
		} `json:"regions"`
	}
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	servers := make([]models.PIAServer, 0, len(data.Regions))
	for _, region := range data.Regions {
		for _, udpServer := range region.Servers.UDP {
			server := models.PIAServer{
				Region:      region.Name,
				ServerName:  udpServer.CN,
				Protocol:    constants.UDP,
				PortForward: region.PortForward,
				IP:          udpServer.IP,
			}
			servers = append(servers, server)
		}
		for _, tcpServer := range region.Servers.TCP {
			server := models.PIAServer{
				Region:      region.Name,
				ServerName:  tcpServer.CN,
				Protocol:    constants.UDP,
				PortForward: region.PortForward,
				IP:          tcpServer.IP,
			}
			servers = append(servers, server)
		}
	}
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			if servers[i].ServerName == servers[j].ServerName {
				return servers[i].Protocol < servers[j].Protocol
			}
			return servers[i].ServerName < servers[j].ServerName
		}
		return servers[i].Region < servers[j].Region
	})

	if u.options.Stdout {
		u.println(stringifyPIAServers(servers))
	}
	u.servers.Pia.Timestamp = u.timeNow().Unix()
	u.servers.Pia.Servers = servers
	return nil
}

func stringifyPIAServers(servers []models.PIAServer) (s string) {
	s = "func PIAServers() []models.PIAServer {\n"
	s += "	return []models.PIAServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
