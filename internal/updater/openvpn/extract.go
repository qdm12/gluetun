package openvpn

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
)

var (
	ErrUnknownProto = errors.New("unknown protocol")
	ErrNoRemoteHost = errors.New("remote host not found")
	ErrNoRemoteIP   = errors.New("remote IP not found")
)

func ExtractProto(b []byte) (tcp, udp bool, err error) {
	lines := strings.Split(string(b), "\n")
	const protoPrefix = "proto "
	for _, line := range lines {
		if !strings.HasPrefix(line, protoPrefix) {
			continue
		}
		s := strings.TrimPrefix(line, protoPrefix)
		s = strings.TrimSpace(s)
		s = strings.ToLower(s)
		switch s {
		case "tcp", "tcp4", "tcp6":
			return true, false, nil
		case "udp", "udp4", "udp6":
			return false, true, nil
		default:
			return false, false, fmt.Errorf("%w: %s", ErrUnknownProto, s)
		}
	}

	// default is UDP if unspecified in openvpn configuration
	return false, true, nil
}

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

func ExtractIPs(b []byte) (ips []net.IP, err error) {
	const rejectIP, rejectDomain = false, true
	ipStrings := extractRemoteHosts(b, rejectIP, rejectDomain)
	if len(ipStrings) == 0 {
		return nil, ErrNoRemoteIP
	}

	sort.Slice(ipStrings, func(i, j int) bool {
		return ipStrings[i] < ipStrings[j]
	})

	ips = make([]net.IP, len(ipStrings))
	for i := range ipStrings {
		ips[i] = net.ParseIP(ipStrings[i])
	}

	return ips, nil
}

func extractRemoteHosts(content []byte, rejectIP, rejectDomain bool) (hosts []string) {
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "remote ") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 1 || fields[1] == "" {
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
