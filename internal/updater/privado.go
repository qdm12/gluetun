package updater

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver"
)

func (u *updater) updatePrivado(ctx context.Context) (err error) {
	servers, warnings, err := findPrivadoServersFromZip(ctx, u.client, u.presolver)
	if u.options.CLI {
		for _, warning := range warnings {
			u.logger.Warn("Privado: %s", warning)
		}
	}
	if err != nil {
		return fmt.Errorf("cannot update Privado servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyPrivadoServers(servers))
	}
	u.servers.Privado.Timestamp = u.timeNow().Unix()
	u.servers.Privado.Servers = servers
	return nil
}

func findPrivadoServersFromZip(ctx context.Context, client *http.Client, presolver resolver.Parallel) (
	servers []models.PrivadoServer, warnings []string, err error) {
	const zipURL = "https://privado.io/apps/ovpn_configs.zip"
	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}

	hosts := make([]string, 0, len(contents))
	for fileName, content := range contents {
		hostname, warning, err := extractHostFromOVPN(content)
		if len(warning) > 0 {
			warnings = append(warnings, warning)
		}
		if err != nil {
			return nil, warnings, fmt.Errorf("%w in %q", err, fileName)
		}
		hosts = append(hosts, hostname)
	}

	const (
		maxFailRatio = 0.1
		maxDuration  = 3 * time.Second
		maxNoNew     = 1
		maxFails     = 2
	)
	settings := resolver.ParallelSettings{
		MaxFailRatio: maxFailRatio,
		Repeat: resolver.RepeatSettings{
			MaxDuration: maxDuration,
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

	for hostname, IPs := range hostToIPs {
		switch len(IPs) {
		case 0:
			warning := fmt.Sprintf("no IP address found for host %q", hostname)
			warnings = append(warnings, warning)
			continue
		case 1:
		default:
			warning := fmt.Sprintf("more than one IP address found for host %q", hostname)
			warnings = append(warnings, warning)
		}
		server := models.PrivadoServer{
			Hostname: hostname,
			IP:       IPs[0],
		}
		servers = append(servers, server)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Hostname < servers[j].Hostname
	})
	return servers, warnings, nil
}

func stringifyPrivadoServers(servers []models.PrivadoServer) (s string) {
	s = "func PrivadoServers() []models.PrivadoServer {\n"
	s += "	return []models.PrivadoServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
