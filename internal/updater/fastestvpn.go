package updater

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateFastestvpn(ctx context.Context) (err error) {
	servers, warnings, err := findFastestvpnServersFromZip(ctx, u.client, u.lookupIP)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("FastestVPN: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update FastestVPN servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyFastestVPNServers(servers))
	}
	u.servers.Fastestvpn.Timestamp = u.timeNow().Unix()
	u.servers.Fastestvpn.Servers = servers
	return nil
}

func findFastestvpnServersFromZip(ctx context.Context, client *http.Client, lookupIP lookupIPFunc) (
	servers []models.FastestvpnServer, warnings []string, err error) {
	const zipURL = "https://support.fastestvpn.com/download/openvpn-tcp-udp-config-files"
	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}

	trailNumberExp := regexp.MustCompile(`[0-9]+$`)

	type Data struct {
		TCP     bool
		UDP     bool
		Country string
	}
	hostToData := make(map[string]Data)

	for fileName, content := range contents {
		const (
			tcpSuffix = "-TCP.ovpn"
			udpSuffix = "-UDP.ovpn"
		)
		var tcp, udp bool
		var suffix string
		switch {
		case strings.HasSuffix(fileName, tcpSuffix):
			suffix = tcpSuffix
			tcp = true
		case strings.HasSuffix(fileName, udpSuffix):
			suffix = udpSuffix
			udp = true
		default:
			warning := `filename "` + fileName + `" does not have a protocol suffix`
			warnings = append(warnings, warning)
			continue
		}

		countryWithNumber := strings.TrimSuffix(fileName, suffix)
		number := trailNumberExp.FindString(countryWithNumber)
		country := countryWithNumber[:len(countryWithNumber)-len(number)]

		host, warning, err := extractHostFromOVPN(content)
		if len(warning) > 0 {
			warnings = append(warnings, warning)
		}
		if err != nil {
			// treat error as warning and go to next file
			warnings = append(warnings, err.Error()+" in "+fileName)
			continue
		}

		data := hostToData[host]
		data.Country = country
		if tcp {
			data.TCP = true
		}
		if udp {
			data.UDP = true
		}
		hostToData[host] = data
	}

	hosts := make([]string, len(hostToData))
	i := 0
	for host := range hostToData {
		hosts[i] = host
		i++
	}

	const repetition = 1
	const timeBetween = 0
	const failOnErr = true
	hostToIPs, _, err := parallelResolve(ctx, lookupIP, hosts, repetition, timeBetween, failOnErr)
	if err != nil {
		return nil, warnings, err
	}

	for host, IPs := range hostToIPs {
		if len(IPs) == 0 {
			warning := fmt.Sprintf("no IP address found for host %q", host)
			warnings = append(warnings, warning)
			continue
		}

		data := hostToData[host]

		server := models.FastestvpnServer{
			Hostname: host,
			TCP:      data.TCP,
			UDP:      data.UDP,
			Country:  data.Country,
			IPs:      uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}

	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Country == servers[j].Country {
			return servers[i].Hostname < servers[j].Hostname
		}
		return servers[i].Country < servers[j].Country
	})

	return servers, warnings, nil
}

func stringifyFastestVPNServers(servers []models.FastestvpnServer) (s string) {
	s = "func FastestvpnServers() []models.FastestvpnServer {\n"
	s += "	return []models.FastestvpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
