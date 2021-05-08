package openvpn

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	ErrNoRemoteHost = errors.New("remote host not found")
	ErrNoRemoteIP   = errors.New("remote IP not found")
)

func ExtractHost(b []byte) (host, warning string, err error) {
	const (
		rejectIP     = true
		rejectDomain = false
	)
	hosts := extractRemoteHosts(b, rejectIP, rejectDomain)
	if len(hosts) == 0 {
		return "", "", ErrNoRemoteHost
	} else if len(hosts) > 1 {
		warning = fmt.Sprintf(
			"only using the first host %q and discarding %d other hosts",
			hosts[0], len(hosts)-1)
	}
	return hosts[0], warning, nil
}

func ExtractIP(b []byte) (ip net.IP, warning string, err error) {
	const (
		rejectIP     = false
		rejectDomain = true
	)
	ips := extractRemoteHosts(b, rejectIP, rejectDomain)
	if len(ips) == 0 {
		return nil, "", ErrNoRemoteIP
	} else if len(ips) > 1 {
		warning = fmt.Sprintf(
			"only using the first IP address %s and discarding %d other hosts",
			ips[0], len(ips)-1)
	}
	return net.ParseIP(ips[0]), warning, nil
}

func extractRemoteHosts(content []byte, rejectIP, rejectDomain bool) (hosts []string) {
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
