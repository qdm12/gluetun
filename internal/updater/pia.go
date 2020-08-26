package updater

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func findPIANewServers() (servers []models.PIAServer, err error) {
	const url = "https://www.privateinternetaccess.com/openvpn/openvpn-ip-nextgen.zip"
	contents, err := fetchAndExtractFiles(url)
	if err != nil {
		return nil, err
	}
	for fileName, content := range contents {
		remoteLines := extractRemoteLinesFromOpenvpn(content)
		if len(remoteLines) == 0 {
			return nil, fmt.Errorf("cannot find any remote lines in %s", fileName)
		}
		IPs := extractIPsFromRemoteLines(remoteLines)
		if len(remoteLines) == 0 {
			return nil, fmt.Errorf("cannot find any IP addresses in %s", fileName)
		}
		region := strings.TrimSuffix(fileName, ".ovpn")
		server := models.PIAServer{
			Region: region,
			IPs:    IPs,
		}
		servers = append(servers, server)
	}
	return servers, nil
}

func stringifyPIAServers(servers []models.PIAServer) (s string) {
	s = "func PIAServers() []models.PIAServer {\n"
	s += "	return []models.PIAServer{\n"
	for _, server := range servers {
		s += server.String() + ",\n"
	}
	s += "	}"
	s += "}"
	return s
}
