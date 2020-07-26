package openvpn

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) Start(ctx context.Context) (stdout io.ReadCloser, waitFn func() error, err error) {
	c.logger.Info("starting openvpn")
	stdout, _, waitFn, err = c.commander.Start(ctx, "openvpn", "--config", string(constants.OpenVPNConf))
	return stdout, waitFn, err
}

func (c *configurator) Version(ctx context.Context) (string, error) {
	output, err := c.commander.Run(ctx, "openvpn", "--version")
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}
	firstLine := strings.Split(output, "\n")[0]
	words := strings.Fields(firstLine)
	if len(words) < 2 {
		return "", fmt.Errorf("openvpn --version: first line is too short: %q", firstLine)
	}
	return words[1], nil
}
