package tinyproxy

import (
	"fmt"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) MakeConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string) error {
	c.logger.Info("%s: generating tinyproxy configuration file", logPrefix)
	lines := generateConf(logLevel, port, user, password)
	return c.fileManager.WriteLinesToFile(string(constants.TinyProxyConf), lines)
}

func generateConf(logLevel models.TinyProxyLogLevel, port uint16, user, password string) (lines []string) {
	confMapping := map[string]string{
		"User":                "nonrootuser",
		"Group":               "tinyproxy",
		"Port":                fmt.Sprintf("%d", port),
		"Timeout":             "600",
		"DefaultErrorFile":    "/usr/share/tinyproxy/default.html",
		"MaxClients":          "100",
		"MinSpareServers":     "5",
		"MaxSpareServers":     "20",
		"StartServers":        "10",
		"MaxRequestsPerChild": "0",
		"DisableViaHeader":    "Yes",
		"LogLevel":            string(logLevel),
		// "StatFile": "/usr/share/tinyproxy/stats.html",
	}
	if len(user) > 0 {
		confMapping["BasicAuth"] = fmt.Sprintf("%s %s", user, password)
	}
	for k, v := range confMapping {
		line := fmt.Sprintf("%s %s", k, v)
		lines = append(lines, line)
	}
	return lines
}
