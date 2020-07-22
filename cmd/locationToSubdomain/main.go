package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/qdm12/golibs/network"
)

func main() {
	os.Exit(_main())
}

func _main() int {
	provider := flag.String("provider", "purevpn", "VPN provider to map location to subdomain, can be 'purevpn'")
	flag.Parse()

	client := network.NewClient(5 * time.Second)
	switch *provider {
	case "purevpn":
		servers, warnings, err := purevpn(client)
		if err != nil {
			fmt.Println(err)
			return 1
		}
		for _, server := range servers {
			fmt.Printf(
				"{subdomain: %q, region: %q, country: %q, city: %q},\n",
				server.subdomain, server.region, server.country, server.city,
			)
		}
		fmt.Print("\n\n")
		for _, warning := range warnings {
			fmt.Println(warning)
		}
	default:
		fmt.Printf("Provider %q is not supported\n", *provider)
		return 1
	}
	return 0
}

type purevpnServer struct {
	region    string
	country   string
	city      string
	subdomain string // without -tcp or -udp suffix
}

func purevpn(client network.Client) (servers []purevpnServer, warnings []string, err error) {
	content, status, err := client.GetContent("https://support.purevpn.com/vpn-servers")
	if err != nil {
		return nil, nil, err
	} else if status != http.StatusOK {
		return nil, nil, fmt.Errorf("HTTP status %d from Purevpn", status)
	}
	const jsonPrefix = "<script>var servers = "
	const jsonSuffix = "</script>"
	s := string(content)
	jsonPrefixIndex := strings.Index(s, jsonPrefix)
	if jsonPrefixIndex == -1 {
		return nil, nil, fmt.Errorf("cannot find prefix %s in html", jsonPrefix)
	}
	if len(s[jsonPrefixIndex:]) == len(jsonPrefix) {
		return nil, nil, fmt.Errorf("no body after json prefix %s", jsonPrefix)
	}
	s = s[jsonPrefixIndex+len(jsonPrefix):]
	endIndex := strings.Index(s, jsonSuffix)
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
	for i := range data {
		if data[i].UDP == "" && data[i].TCP == "" {
			warnings = append(warnings, fmt.Sprintf("server %s %s %s does not support TCP and UDP for openvpn", data[i].Region, data[i].Country, data[i].City))
			continue
		}
		if data[i].UDP == "" || data[i].TCP == "" {
			warnings = append(warnings, fmt.Sprintf("server %s %s %s does not support TCP or udp for openvpn", data[i].Region, data[i].Country, data[i].City))
		}
		servers = append(servers, purevpnServer{
			region:    data[i].Region,
			country:   data[i].Country,
			city:      data[i].City,
			subdomain: strings.TrimSuffix(data[i].TCP, "-tcp.pointtoserver.com"),
		})
	}
	return servers, warnings, nil
}
