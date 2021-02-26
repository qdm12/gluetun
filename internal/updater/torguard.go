package updater

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateTorguard(ctx context.Context) (err error) {
	servers, warnings, err := findTorguardServersFromZip(ctx, u.client)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Torguard: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Torguard servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyTorguardServers(servers))
	}
	u.servers.Torguard.Timestamp = u.timeNow().Unix()
	u.servers.Torguard.Servers = servers
	return nil
}

func findTorguardServersFromZip(ctx context.Context, client *http.Client) (
	servers []models.TorguardServer, warnings []string, err error) {
	// Note: all servers do both TCP and UDP
	const zipURL = "https://torguard.net/downloads/OpenVPN-TCP-Linux.zip"

	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}

	for fileName, content := range contents {
		var server models.TorguardServer

		const prefix = "TorGuard."
		const suffix = ".ovpn"
		s := strings.TrimPrefix(fileName, prefix)
		s = strings.TrimSuffix(s, suffix)

		switch {
		case strings.Count(s, ".") == 1 && !strings.HasPrefix(s, "USA"):
			parts := strings.Split(s, ".")
			server.Country = parts[0]
			server.City = parts[1]
		case strings.HasPrefix(s, "USA"):
			server.Country = "USA"
			s = strings.TrimPrefix(s, "USA-")
			s = strings.ReplaceAll(s, "-", " ")
			s = strings.ReplaceAll(s, ".", " ")
			s = strings.ToLower(s)
			s = strings.Title(s)
			server.City = s
		default:
			server.Country = s
		}

		hostnames := extractRemoteHostsFromOpenvpn(content, true, false)
		if len(hostnames) != 1 {
			warning := "found " + strconv.Itoa(len(hostnames)) +
				" hostname(s) instead of 1 in " + fileName
			warnings = append(warnings, warning)
			continue
		}
		server.Hostname = hostnames[0]

		IPs := extractRemoteHostsFromOpenvpn(content, false, true)
		if len(IPs) != 1 {
			warning := "found " + strconv.Itoa(len(IPs)) +
				" IP(s) instead of 1 in " + fileName
			warnings = append(warnings, warning)
			continue
		}
		server.IP = net.ParseIP(IPs[0])
		if server.IP == nil {
			warnings = append(warnings, "IP address "+IPs[0]+" is not valid in file "+fileName)
		}

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

func stringifyTorguardServers(servers []models.TorguardServer) (s string) {
	s = "func TorguardServers() []models.TorguardServer {\n"
	s += "	return []models.TorguardServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
