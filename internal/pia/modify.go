package pia

import (
	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func modifyLines(lines, IPs []string, port string) (modifiedLines []string, err error) {
	// Remove lines
	for _, line := range lines {
		if strings.Contains(line, "privateinternetaccess.com") ||
			strings.Contains(line, "resolve-retry") {
			continue
		}
		modifiedLines = append(modifiedLines, line)
	}
	// Add lines
	for _, IP := range IPs {
		modifiedLines = append(modifiedLines, fmt.Sprintf("remote %s %s", IP, port))
	}
	modifiedLines = append(modifiedLines, "auth-user-pass "+constants.OpenVPNAuthConf)
	modifiedLines = append(modifiedLines, "auth-retry nointeract")
	modifiedLines = append(modifiedLines, "pull-filter ignore \"auth-token\"") // prevent auth failed loops
	modifiedLines = append(modifiedLines, "user nonrootuser")
	modifiedLines = append(modifiedLines, "mute-replay-warnings")
	return modifiedLines, nil
}
