package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	provider := flag.String("provider", "nordvpn", "VPN provider to map region to IP addresses using their API, can be 'nordvpn'")
	flag.Parse()

	client := network.NewClient(30 * time.Second) // big file so 30 seconds
	switch *provider {
	case "nordvpn":
		servers, ignoredServers, err := nordvpn(client)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		for _, server := range servers {
			fmt.Printf(
				"{Region: %q, Number: %d, TCP: %t, UDP: %t, IP: net.IP{%s}},\n",
				server.Region, server.Number, server.TCP, server.UDP, strings.ReplaceAll(server.IP.String(), ".", ", "),
			)
		}
		fmt.Print("\n\n")
		for _, serverName := range ignoredServers {
			fmt.Printf("ignored server %q because it does not support both UDP and TCP\n", serverName)
		}
	default:
		fmt.Printf("Provider %q is not supported\n", *provider)
		return 1
	}
	return 0
}

func nordvpn(client network.Client) (servers []models.NordvpnServer, ignoredServers []string, err error) {
	content, status, err := client.GetContent("https://nordvpn.com/api/server")
	if err != nil {
		return nil, nil, err
	} else if status != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP status %d from NordVPN API", status)
	}
	response := []struct {
		IPAddress string `json:"ip_address"`
		Name      string `json:"name"`
		Country   string `json:"country"`
		Features  struct {
			UDP bool `json:"openvpn_udp"`
			TCP bool `json:"openvpn_tcp"`
		} `json:"features"`
	}{}
	if err := json.Unmarshal(content, &response); err != nil {
		return nil, nil, err
	}

	for _, element := range response {
		if !element.Features.TCP && !element.Features.UDP {
			ignoredServers = append(ignoredServers, element.Name)
		}
		ip := net.ParseIP(element.IPAddress)
		if ip == nil {
			return nil, nil, fmt.Errorf("IP address %q is not valid for server %q", element.IPAddress, element.Name)
		}
		i := strings.IndexRune(element.Name, '#')
		if i < 0 {
			return nil, nil, fmt.Errorf("No ID in server name %q", element.Name)
		}
		idString := element.Name[i+1:]
		idUint64, err := strconv.ParseUint(idString, 10, 16)
		if err != nil {
			return nil, nil, fmt.Errorf("Bad ID in server name %q", element.Name)
		}
		id := uint16(idUint64)
		server := models.NordvpnServer{
			Region: element.Country,
			Number: id,
			IP:     ip,
			TCP:    element.Features.TCP,
			UDP:    element.Features.UDP,
		}
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Region == servers[j].Region {
			return servers[i].Number < servers[j].Number
		}
		return servers[i].Region < servers[j].Region
	})
	return servers, ignoredServers, nil
}
