package updater

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateVyprvpn(ctx context.Context) (err error) {
	servers, err := findVyprvpnServers(ctx, u.lookupIP)
	if err != nil {
		return fmt.Errorf("cannot update Vyprvpn servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyVyprvpnServers(servers))
	}
	u.servers.Vyprvpn.Timestamp = u.timeNow().Unix()
	u.servers.Vyprvpn.Servers = servers
	return nil
}

func findVyprvpnServers(ctx context.Context, lookupIP lookupIPFunc) (servers []models.VyprvpnServer, err error) {
	const zipURL = "https://support.vyprvpn.com/hc/article_attachments/360052617332/Vypr_OpenVPN_20200320.zip"
	contents, err := fetchAndExtractFiles(zipURL)
	if err != nil {
		return nil, err
	}
	for fileName, content := range contents {
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return nil, fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		hosts := extractHostnamesFromRemoteLines(remoteLines)
		if len(hosts) == 0 {
			return nil, fmt.Errorf("cannot find any hosts in %s", fileName)
		}
		var IPs []net.IP
		for _, host := range hosts {
			newIPs, err := lookupIP(ctx, host)
			if err != nil {
				return nil, err
			}
			IPs = append(IPs, newIPs...)
		}
		region := strings.TrimSuffix(fileName, ".ovpn")
		region = strings.ReplaceAll(region, " - ", " ")
		server := models.VyprvpnServer{
			Region: region,
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	return servers, nil
}

func stringifyVyprvpnServers(servers []models.VyprvpnServer) (s string) {
	s = "func VyprvpnServers() []models.VyprvpnServer {\n"
	s += "	return []models.VyprvpnServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
