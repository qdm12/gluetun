package updater

import (
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

func extractHostnamesFromRemoteLines(remoteLines []string) (hostnames []string) {
	for _, remoteLine := range remoteLines {
		fields := strings.Fields(remoteLine)
		if len(fields[1]) == 0 {
			continue
		}
		hostnames = append(hostnames, fields[1])
	}
	return hostnames
}
