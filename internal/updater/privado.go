package updater

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

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
		hosts := extractHostnamesFromRemoteLines(remoteLines)
		if len(hosts) == 0 {
			return nil, warnings, fmt.Errorf("cannot find any hosts in %s", fileName)
		} else if len(hosts) > 1 {
			warning := fmt.Sprintf("more than one host in %q, only taking first one %q into account", fileName, hosts[0])
			warnings = append(warnings, warning)
		}
		host := hosts[0]
		if net.ParseIP(host) != nil {
			warning := fmt.Sprintf("ignoring IP address host %q in %s", host, fileName)
			warnings = append(warnings, warning)
			continue
		}
		const repetition = 1
		IPs, err := resolveRepeat(ctx, lookupIP, host, repetition)
		if err != nil {
			return nil, warnings, err
		} else if len(IPs) == 0 {
			warning := fmt.Sprintf("no IP address found for host %q", host)
			warnings = append(warnings, warning)
			continue
		}
		subdomain := strings.TrimSuffix(host, ".vpn.privado.io")
		parts := strings.Split(subdomain, "-")
		const expectedParts = 2
		if len(parts) != expectedParts {
			warning := fmt.Sprintf("malformed subdomain %q: cannot find city and server number", subdomain)
			warnings = append(warnings, warning)
			continue
		}
		city, serverNumberString := parts[0], parts[1]
		serverNumberInt, err := strconv.ParseInt(serverNumberString, 10, 16)
		if err != nil {
			warning := fmt.Sprintf("malformed server number %q: %s", serverNumberString, err)
			warnings = append(warnings, warning)
			continue
		}
		server := models.PrivadoServer{
			City:   city,
			Number: uint16(serverNumberInt),
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}

	sort.Slice(servers, func(i, j int) bool {
		keyA := servers[i].City + fmt.Sprintf("%d", servers[i].Number)
		keyB := servers[j].City + fmt.Sprintf("%d", servers[j].Number)
		return keyA < keyB
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
