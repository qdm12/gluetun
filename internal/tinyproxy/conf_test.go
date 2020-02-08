package tinyproxy

import (
	"testing"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_generateConf(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		logLevel models.TinyProxyLogLevel
		port     uint16
		user     string
		password string
		lines    []string
	}{
		"No credentials": {
			logLevel: constants.TinyProxyInfoLevel,
			port:     2000,
			lines: []string{
				"DefaultErrorFile \"/usr/share/tinyproxy/default.html\"",
				"DisableViaHeader Yes",
				"Group tinyproxy",
				"LogLevel Info",
				"MaxClients 100",
				"MaxRequestsPerChild 0",
				"MaxSpareServers 20",
				"MinSpareServers 5",
				"Port 2000",
				"StartServers 10",
				"Timeout 600",
				"User nonrootuser",
			},
		},
		"With credentials": {
			logLevel: constants.TinyProxyErrorLevel,
			port:     2000,
			user:     "abc",
			password: "def",
			lines: []string{
				"BasicAuth abc def",
				"DefaultErrorFile \"/usr/share/tinyproxy/default.html\"",
				"DisableViaHeader Yes",
				"Group tinyproxy",
				"LogLevel Error",
				"MaxClients 100",
				"MaxRequestsPerChild 0",
				"MaxSpareServers 20",
				"MinSpareServers 5",
				"Port 2000",
				"StartServers 10",
				"Timeout 600",
				"User nonrootuser",
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			lines := generateConf(tc.logLevel, tc.port, tc.user, tc.password)
			assert.Equal(t, tc.lines, lines)
		})
	}
}
