package pia

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/golibs/verification"
)

func parseConfig(lines []string, verifier verification.Verifier, lookupIP func(host string) ([]net.IP, error)) (
	IPs []string, port, device string, err error) {
	remoteLineFound := false
	deviceLineFound := false
	for _, line := range lines {
		if strings.HasPrefix(line, "remote ") {
			remoteLineFound = true
			words := strings.Split(line, " ")
			if len(words) != 3 {
				return nil, "", "", fmt.Errorf("line %q misses information", line)
			}
			host := words[1]
			if err := verifier.VerifyPort(words[2]); err != nil {
				return nil, "", "", fmt.Errorf("line %q has an invalid port: %w", line, err)
			}
			port = words[2]
			netIPs, err := lookupIP(host) // TODO use Unbound
			if err != nil {
				return nil, "", "", err
			}
			for _, netIP := range netIPs {
				IPs = append(IPs, netIP.String())
			}
		} else if strings.HasPrefix(line, "dev ") {
			deviceLineFound = true
			words := strings.Split(line, " ")
			if len(words) != 2 {
				return nil, "", "", fmt.Errorf("line %q misses information", line)
			}
			device = words[1]
			if device != "tun" && device != "tap" {
				return nil, "", "", fmt.Errorf("device %q is not valid", device)
			}
		}
	}
	if remoteLineFound && deviceLineFound {
		return IPs, port, device, nil
	} else if !remoteLineFound {
		return nil, "", "", fmt.Errorf("remote line not found in Openvpn configuration")
	} else {
		return nil, "", "", fmt.Errorf("device line not found in Openvpn configuration")
	}
}
