package updater

import (
	"fmt"
	"sort"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updatePIA() (err error) {
	const zipURL = "https://www.privateinternetaccess.com/openvpn/openvpn-ip-nextgen.zip"
	servers, err := findPIAServersFromURL(zipURL)
	if err != nil {
		return fmt.Errorf("cannot update PIA servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyPIAServers(servers))
	}
	u.servers.Pia.Timestamp = u.timeNow().Unix()
	u.servers.Pia.Servers = servers
	return nil
}

func (u *updater) updatePIAOld() (err error) {
	const zipURL = "https://www.privateinternetaccess.com/openvpn/openvpn-ip.zip"
	servers, err := findPIAServersFromURL(zipURL)
	if err != nil {
		return fmt.Errorf("cannot update old PIA servers: %w", err)
	}
	if u.options.Stdout {
		u.println(stringifyPIAOldServers(servers))
	}
	u.servers.PiaOld.Timestamp = u.timeNow().Unix()
	u.servers.PiaOld.Servers = servers
	return nil
}

func findPIAServersFromURL(zipURL string) (servers []models.PIAServer, err error) {
	contents, err := fetchAndExtractFiles(zipURL)
	if err != nil {
		return nil, err
	}
	for fileName, content := range contents {
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return nil, fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		IPs := extractIPsFromRemoteLines(remoteLines)
		if len(IPs) == 0 {
			return nil, fmt.Errorf("cannot find any IP addresses in %s", fileName)
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
	return servers, nil
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
