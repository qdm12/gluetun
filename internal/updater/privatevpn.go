package updater

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePrivatevpn(ctx context.Context) (err error) {
	servers, warnings, err := findPrivatevpnServersFromZip(ctx, u.client, u.lookupIP)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Privatevpn: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Privatevpn servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyPrivatevpnServers(servers))
	}
	u.servers.Privatevpn.Timestamp = u.timeNow().Unix()
	u.servers.Privatevpn.Servers = servers
	return nil
}

func findPrivatevpnServersFromZip(ctx context.Context, client *http.Client, lookupIP lookupIPFunc) (
	servers []models.PrivatevpnServer, warnings []string, err error) {
	// Note: all servers do both TCP and UDP
	const zipURL = "https://privatevpn.com/client/PrivateVPN-TUN.zip"

	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}

	trailingNumber := regexp.MustCompile(` [0-9]+$`)
	countryCodes := constants.CountryCodes()

	uniqueServers := map[string]models.PrivatevpnServer{} // key is the hostname

	for fileName, content := range contents {
		const prefix = "PrivateVPN-"
		const suffix = "-TUN-443.ovpn"

		if !strings.HasSuffix(fileName, suffix) {
			continue // only process TCP servers as they're the same
		}

		var server models.PrivatevpnServer

		s := strings.TrimPrefix(fileName, prefix)
		s = strings.TrimSuffix(s, suffix)
		s = trailingNumber.ReplaceAllString(s, "")

		parts := strings.Split(s, "-")
		var countryCode string
		countryCode, server.City = parts[0], parts[1]
		countryCode = strings.ToLower(countryCode)
		var countryCodeOK bool
		server.Country, countryCodeOK = countryCodes[countryCode]
		if !countryCodeOK {
			warnings = append(warnings, "unknown country code: "+countryCode)
			server.Country = countryCode
		}

		var warning string
		server.Hostname, warning, err = extractHostFromOVPN(content)
		if len(warning) > 0 {
			warnings = append(warnings, warning)
		}
		if err != nil {
			return nil, warnings, err
		}
		if len(warning) > 0 {
			continue
		}

		uniqueServers[server.Hostname] = server
	}

	hostnames := make([]string, len(uniqueServers))
	i := 0
	for hostname := range uniqueServers {
		hostnames[i] = hostname
		i++
	}

	const failOnError = false
	hostToIPs, newWarnings, _ := parallelResolve(ctx, lookupIP, hostnames, 5, time.Second, failOnError)
	if len(newWarnings) > 0 {
		warnings = append(warnings, newWarnings...)
	}

	for hostname, server := range uniqueServers {
		ips := hostToIPs[hostname]
		if len(ips) == 0 {
			continue
		}
		server.IPs = ips
		servers = append(servers, server)
	}

	sort.Slice(servers, func(i, j int) bool {
		if servers[i].Country == servers[j].Country {
			if servers[i].City == servers[j].City {
				return servers[i].Hostname < servers[j].Hostname
			}
			return servers[i].City < servers[j].City
		}
		return servers[i].Country < servers[j].Country
	})

	return servers, warnings, nil
}

func stringifyPrivatevpnServers(servers []models.PrivatevpnServer) (s string) {
	s = "func PrivatevpnServers() []models.PrivatevpnServer {\n"
	s += "	return []models.PrivatevpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
