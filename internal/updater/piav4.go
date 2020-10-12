package updater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePIA() (err error) {
	const url = "https://serverlist.piaservers.net/vpninfo/servers/v4"
	response, err := u.httpGet(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	} else if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: %s", response.Status, strings.ReplaceAll(string(b), "\n", ""))
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
		server := models.PIAServer{
			Region:      region.Name,
			PortForward: region.PortForward,
		}
		for _, udpServer := range region.Servers.UDP {
			if len(server.OpenvpnUDP.CN) > 0 && server.OpenvpnUDP.CN != udpServer.CN {
				return fmt.Errorf("CN is different for UDP for region %q: %q and %q", region.Name, server.OpenvpnUDP.CN, udpServer.CN)
			}
			if udpServer.IP != nil {
				server.OpenvpnUDP.IPs = append(server.OpenvpnUDP.IPs, udpServer.IP)
			}
		}
		for _, tcpServer := range region.Servers.TCP {
			if len(server.OpenvpnTCP.CN) > 0 && server.OpenvpnTCP.CN != tcpServer.CN {
				return fmt.Errorf("CN is different for TCP for region %q: %q and %q", region.Name, server.OpenvpnTCP.CN, tcpServer.CN)
			}
			if tcpServer.IP != nil {
				server.OpenvpnTCP.IPs = append(server.OpenvpnTCP.IPs, tcpServer.IP)
			}
		}
		if server.OpenvpnTCP.CN != server.OpenvpnUDP.CN {
			return fmt.Errorf("not the same: %q, %q", server.OpenvpnTCP.CN, server.OpenvpnUDP.CN)
		}
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
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
