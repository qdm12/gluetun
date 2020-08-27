package updater

import (
	"net"
	"strings"
)

func extractRemoteLinesFromOpenvpn(content []byte) (remoteLines []string) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "remote ") {
			remoteLines = append(remoteLines, line)
		}
	}
	return remoteLines
}

func extractIPsFromRemoteLines(remoteLines []string) (ips []net.IP) {
	for _, remoteLine := range remoteLines {
		fields := strings.Fields(remoteLine)
		ip := net.ParseIP(fields[1])
		if ip == nil { // not an IP address
			continue
		}
		ips = append(ips, ip)
	}
	return ips
}
