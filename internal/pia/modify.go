package pia

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) ModifyLines(lines []string, IPs []net.IP, port uint16) (modifiedLines []string) {
	c.logger.Info("%s: adapting openvpn configuration for server IP addresses and port %d", logPrefix, port)
	// Remove lines
	for _, line := range lines {
		if strings.Contains(line, "privateinternetaccess.com") ||
			strings.Contains(line, "resolv-retry") {
			continue
		}
		modifiedLines = append(modifiedLines, line)
	}
	// Add lines
	for _, IP := range IPs {
		modifiedLines = append(modifiedLines, fmt.Sprintf("remote %s %d", IP.String(), port))
	}
	modifiedLines = append(modifiedLines, "auth-user-pass "+string(constants.OpenVPNAuthConf))
	modifiedLines = append(modifiedLines, "auth-retry nointeract")
	modifiedLines = append(modifiedLines, "pull-filter ignore \"auth-token\"") // prevent auth failed loops
	modifiedLines = append(modifiedLines, "user nonrootuser")
	modifiedLines = append(modifiedLines, "mute-replay-warnings")
	return modifiedLines
}
