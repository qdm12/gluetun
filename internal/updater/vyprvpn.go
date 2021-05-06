package updater

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func (u *updater) updateVyprvpn(ctx context.Context) (err error) {
	servers, warnings, err := findVyprvpnServers(ctx, u.client, u.presolver)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Vyprvpn: %s", warning)
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

func findVyprvpnServers(ctx context.Context, client *http.Client, presolver resolver.Parallel) (
	servers []models.VyprvpnServer, warnings []string, err error) {
	const zipURL = "https://support.vyprvpn.com/hc/article_attachments/360052617332/Vypr_OpenVPN_20200320.zip"
	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}

	hostToRegion := make(map[string]string, len(contents))
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
		region := strings.TrimSuffix(fileName, ".ovpn")
		region = strings.ReplaceAll(region, " - ", " ")
		hostToRegion[host] = region
	}

	hosts := make([]string, len(hostToRegion))
	i := 0
	for host := range hostToRegion {
		hosts[i] = host
		i++
	}

	const (
		maxFailRatio = 0.1
		maxNoNew     = 2
		maxFails     = 2
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration: time.Second,
			MaxNoNew:    maxNoNew,
			MaxFails:    maxFails,
			SortIPs:     true,
		},
	}
	hostToIPs, newWarnings, err := presolver.Resolve(ctx, hosts, settings)
	warnings = append(warnings, newWarnings...)
	if err != nil {
		return nil, warnings, err
	}

	servers = make([]models.VyprvpnServer, 0, len(hostToIPs))
	for host, IPs := range hostToIPs {
		server := models.VyprvpnServer{
			Region: hostToRegion[host],
			IPs:    IPs,
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
