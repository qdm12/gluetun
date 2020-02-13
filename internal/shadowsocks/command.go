package shadowsocks

import (
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) Start(server string, port uint16, password string, log bool) (stdout io.ReadCloser, waitFn func() error, err error) {
	c.logger.Info("%s: starting shadowsocks server", logPrefix)
	args := []string{
		"-c", string(constants.ShadowsocksConf),
		"-p", fmt.Sprintf("%d", port),
		"-k", password,
	}
	if log {
		args = append(args, "-v")
	}
	stdout, _, waitFn, err = c.commander.Start("ss-server", args...)
	return stdout, waitFn, err
}

// Version obtains the version of the installed shadowsocks server
func (c *configurator) Version() (string, error) {
	output, err := c.commander.Run("ss-server", "-h")
	if err != nil {
		return "", err
	}
	lines := strings.Split(output, "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("ss-server -h: not enough lines in %q", output)
	}
	words := strings.Fields(lines[1])
	if len(words) < 2 {
		return "", fmt.Errorf("ss-server -h: line 2 is too short: %q", lines[1])
	}
	return words[1], nil
}
