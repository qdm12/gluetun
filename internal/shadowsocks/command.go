package shadowsocks

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func (c *configurator) Start(ctx context.Context, server string, port uint16, password string, log bool) (stdout, stderr io.ReadCloser, waitFn func() error, err error) {
	c.logger.Info("starting shadowsocks server")
	args := []string{
		"-c", string(constants.ShadowsocksConf),
		"-p", fmt.Sprintf("%d", port),
		"-k", password,
	}
	if log {
		args = append(args, "-v")
	}
	stdout, stderr, waitFn, err = c.commander.Start(ctx, "ss-server", args...)
	return stdout, stderr, waitFn, err
}

// Version obtains the version of the installed shadowsocks server
func (c *configurator) Version(ctx context.Context) (string, error) {
	output, err := c.commander.Run(ctx, "ss-server", "-h")
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
