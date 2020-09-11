package updater

import (
	"context"
	"fmt"
	"net"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePIA() (err error) {
	const zipURL = "https://www.privateinternetaccess.com/openvpn/openvpn-ip-nextgen.zip"
	contents, err := fetchAndExtractFiles(zipURL)
	if err != nil {
		return err
	}
	servers := make([]models.PIAServer, 0, len(contents))
	for fileName, content := range contents {
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		IPs := extractIPsFromRemoteLines(remoteLines)
		if len(IPs) == 0 {
			return fmt.Errorf("cannot find any IP addresses in %s", fileName)
		}
		region := strings.TrimSuffix(fileName, ".ovpn")
		server := models.PIAServer{
			Region: region,
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	if u.options.Stdout {
		u.println(stringifyPIAServers(servers))
	}
	u.servers.Pia.Timestamp = u.timeNow().Unix()
	u.servers.Pia.Servers = servers
	return nil
}

func (u *updater) updatePIAOld(ctx context.Context) (err error) {
	const zipURL = "https://www.privateinternetaccess.com/openvpn/openvpn.zip"
	contents, err := fetchAndExtractFiles(zipURL)
	if err != nil {
		return err
	}
	servers := make([]models.PIAServer, 0, len(contents))
	for fileName, content := range contents {
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		hosts := extractHostnamesFromRemoteLines(remoteLines)
		if len(hosts) == 0 {
			return fmt.Errorf("cannot find any hosts in %s", fileName)
		}
		var IPs []net.IP
		for _, host := range hosts {
			newIPs, err := resolveRepeat(ctx, u.lookupIP, host, 3)
			if err != nil {
				return err
			}
			IPs = append(IPs, newIPs...)
		}
		region := strings.TrimSuffix(fileName, ".ovpn")
		server := models.PIAServer{
			Region: region,
			IPs:    uniqueSortedIPs(IPs),
		}
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	if u.options.Stdout {
		u.println(stringifyPIAOldServers(servers))
	}
	u.servers.PiaOld.Timestamp = u.timeNow().Unix()
	u.servers.PiaOld.Servers = servers
	return nil
}

func stringifyPIAServers(servers []models.PIAServer) (s string) {
	s = "func PIAServers() []models.PIAServer {\n"
	s += "	return []models.PIAServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}

func stringifyPIAOldServers(servers []models.PIAServer) (s string) {
	s = "func PIAOldServers() []models.PIAServer {\n"
	s += "	return []models.PIAServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}
