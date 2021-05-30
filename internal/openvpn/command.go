package openvpn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

func (c *configurator) Start(ctx context.Context) (
	stdoutLines, stderrLines chan string, waitError chan error, err error) {
	c.logger.Info("starting openvpn")
	return c.commander.Start(ctx, "openvpn", "--config", constants.OpenVPNConf)
}

var ErrVersionTooShort = errors.New("version output is too short")

func (c *configurator) Version(ctx context.Context) (string, error) {
	output, err := c.commander.Run(ctx, "openvpn", "--version")
	if err != nil && err.Error() != "exit status 1" {
		return "", err
	}
	firstLine := strings.Split(output, "\n")[0]
	words := strings.Fields(firstLine)
	const minWords = 2
	if len(words) < minWords {
		return "", fmt.Errorf("%w: %s", ErrVersionTooShort, firstLine)
	}
	return words[1], nil
}
