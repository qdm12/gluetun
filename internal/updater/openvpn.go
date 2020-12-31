package updater

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	errRemoteHostNotFound = errors.New("remote host not found")
)

func extractHostFromOVPN(b []byte) (host, warning string, err error) {
	const (
		rejectIP     = true
		rejectDomain = false
	)
	hosts := extractRemoteHostsFromOpenvpn(b, rejectIP, rejectDomain)
	if len(hosts) == 0 {
		return "", "", errRemoteHostNotFound
	} else if len(hosts) > 1 {
		warning = fmt.Sprintf(
			"only using the first host %q and discarding %d other hosts",
			hosts[0], len(hosts)-1)
	}
	return hosts[0], warning, nil
}

func extractRemoteHostsFromOpenvpn(content []byte,
	rejectIP, rejectDomain bool) (hosts []string) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if !strings.HasPrefix(line, "remote ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 1 || len(fields[1]) == 0 {
			continue
		}
		host := fields[1]
		parsedIP := net.ParseIP(host)
		if (rejectIP && parsedIP != nil) ||
			(rejectDomain && parsedIP == nil) {
			continue
		}
		hosts = append(hosts, host)
	}
	return hosts
}
