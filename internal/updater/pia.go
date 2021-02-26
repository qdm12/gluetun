package updater

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePIA(ctx context.Context) (err error) {
	const url = "https://serverlist.piaservers.net/vpninfo/servers/v5"

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := u.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s for %s", ErrHTTPStatusCodeNotOK, response.Status, url)
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if err := response.Body.Close(); err != nil {
		return err
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

	// Deduplicate servers with the same IP address
	type NetProtocols struct {
		tcp, udp bool
	}
	ipToProtocols := make(map[string]NetProtocols)

	for _, region := range data.Regions {
		for _, udpServer := range region.Servers.UDP {
			protocols := ipToProtocols[udpServer.IP.String()]
			protocols.udp = true
			ipToProtocols[udpServer.IP.String()] = protocols
		}
		for _, tcpServer := range region.Servers.TCP {
			protocols := ipToProtocols[tcpServer.IP.String()]
			protocols.tcp = true
			ipToProtocols[tcpServer.IP.String()] = protocols
		}
	}

	servers := make([]models.PIAServer, 0, len(ipToProtocols)) // set the capacity, not the length of the slice
	for _, region := range data.Regions {
		for _, udpServer := range region.Servers.UDP {
			protocols, ok := ipToProtocols[udpServer.IP.String()]
			if !ok { // already added that IP for a server
				continue
			}
			server := models.PIAServer{
				Region:      region.Name,
				ServerName:  udpServer.CN,
				TCP:         protocols.tcp,
				UDP:         protocols.udp,
				PortForward: region.PortForward,
				IP:          udpServer.IP,
			}
			delete(ipToProtocols, udpServer.IP.String())
			servers = append(servers, server)
		}
		for _, tcpServer := range region.Servers.TCP {
			protocols, ok := ipToProtocols[tcpServer.IP.String()]
			if !ok { // already added that IP for a server
				continue
			}
			server := models.PIAServer{
				Region:      region.Name,
				ServerName:  tcpServer.CN,
				TCP:         protocols.tcp,
				UDP:         protocols.udp,
				PortForward: region.PortForward,
				IP:          tcpServer.IP,
			}
			delete(ipToProtocols, tcpServer.IP.String())
			servers = append(servers, server)
		}
	}
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
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
