package updater

import (
	"context"
	"fmt"
	"net"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/network"
)

func (u *updater) updatePrivado(ctx context.Context) (err error) {
	servers, warnings, err := findPrivadoServersFromZip(ctx, u.client, u.lookupIP)
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

func findPrivadoServersFromZip(ctx context.Context, client network.Client, lookupIP lookupIPFunc) (
	servers []models.PrivadoServer, warnings []string, err error) {
	const zipURL = "https://privado.io/apps/ovpn_configs.zip"
	contents, err := fetchAndExtractFiles(ctx, client, zipURL)
	if err != nil {
		return nil, nil, err
	}
	for fileName, content := range contents {
		if err := ctx.Err(); err != nil {
			return nil, warnings, err
		}
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return nil, warnings, fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		hostnames := extractHostnamesFromRemoteLines(remoteLines)
		if len(hostnames) == 0 {
			return nil, warnings, fmt.Errorf("cannot find any hosts in %s", fileName)
		} else if len(hostnames) > 1 {
			warning := fmt.Sprintf("more than one host in %q, only taking first one %q into account", fileName, hostnames[0])
			warnings = append(warnings, warning)
		}
		hostname := hostnames[0]
		if net.ParseIP(hostname) != nil {
			warning := fmt.Sprintf("ignoring IP address host %q in %s", hostname, fileName)
			warnings = append(warnings, warning)
			continue
		}
		const repetition = 1
		IPs, err := resolveRepeat(ctx, lookupIP, hostname, repetition)
		switch {
		case err != nil:
			return nil, warnings, err
		case len(IPs) == 0:
			warning := fmt.Sprintf("no IP address found for host %q", hostname)
			warnings = append(warnings, warning)
			continue
		case len(IPs) > 1:
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
