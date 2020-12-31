package updater

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

func (u *updater) updateVyprvpn(ctx context.Context) (err error) {
	servers, warnings, err := findVyprvpnServers(ctx, u.client, u.lookupIP)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Privado: %s", warning)
		}
	}
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

func findVyprvpnServers(ctx context.Context, client network.Client, lookupIP lookupIPFunc) (
	servers []models.VyprvpnServer, warnings []string, err error) {
	const zipURL = "https://support.vyprvpn.com/hc/article_attachments/360052617332/Vypr_OpenVPN_20200320.zip"
	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}
	for fileName, content := range contents {
		if err := ctx.Err(); err != nil {
			return nil, warnings, err
		}
		host, warning, err := extractHostFromOVPN(content)
		if len(warning) > 0 {
			warnings = append(warnings, warning)
		}
		if err != nil {
			return nil, warnings, fmt.Errorf("%w in %s", err, fileName)
		}
		const repetitions = 1
		IPs, err := resolveRepeat(ctx, lookupIP, host, repetitions)
		if err != nil {
			return nil, warnings, err
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
	return servers, warnings, nil
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
