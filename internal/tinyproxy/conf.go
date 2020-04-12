package tinyproxy

import (
	"fmt"
	"sort"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) MakeConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string, uid, gid int) error {
	c.logger.Info("generating tinyproxy configuration file")
	lines := generateConf(logLevel, port, user, password, uid, gid)
	return c.fileManager.WriteLinesToFile(string(constants.TinyProxyConf),
		lines,
		files.Ownership(uid, gid),
		files.Permissions(0400))
}

func generateConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string, uid, gid int) (lines []string) {
	confMapping := map[string]string{
		"User":                fmt.Sprintf("%d", uid),
		"Group":               fmt.Sprintf("%d", gid),
		"Port":                fmt.Sprintf("%d", port),
		"Timeout":             "600",
		"DefaultErrorFile":    "\"/usr/share/tinyproxy/default.html\"",
		"MaxClients":          "100",
		"MinSpareServers":     "5",
		"MaxSpareServers":     "20",
		"StartServers":        "10",
		"MaxRequestsPerChild": "0",
		"DisableViaHeader":    "Yes",
		"LogLevel":            string(logLevel),
		// "StatFile": "\"/usr/share/tinyproxy/stats.html\"",
	}
	if len(user) > 0 {
		confMapping["BasicAuth"] = fmt.Sprintf("%s %s", user, password)
	}
	for k, v := range confMapping {
		line := fmt.Sprintf("%s %s", k, v)
		lines = append(lines, line)
	}
	sort.Slice(lines, func(i, j int) bool {
		return lines[i] < lines[j]
	})
	return lines
}
