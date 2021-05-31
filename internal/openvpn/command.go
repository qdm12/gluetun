package openvpn

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
)

var ErrVersionUnknown = errors.New("OpenVPN version is unknown")

const (
	binOpenvpn24 = "openvpn2.4"
	binOpenvpn25 = "openvpn"
)

func (c *configurator) Start(ctx context.Context, version string) (
	stdoutLines, stderrLines chan string, waitError chan error, err error) {
	var bin string
	switch version {
	case constants.Openvpn24:
		bin = binOpenvpn24
	case constants.Openvpn25:
		bin = binOpenvpn25
	default:
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrVersionUnknown, version)
	}

	c.logger.Info("starting OpenVPN " + version)

	return c.commander.Start(ctx, bin, "--config", constants.OpenVPNConf)
}

func (c *configurator) Version24(ctx context.Context) (version string, err error) {
	return c.version(ctx, binOpenvpn24)
}

func (c *configurator) Version25(ctx context.Context) (version string, err error) {
	return c.version(ctx, binOpenvpn25)
}

var ErrVersionTooShort = errors.New("version output is too short")

func (c *configurator) version(ctx context.Context, binName string) (version string, err error) {
	output, err := c.commander.Run(ctx, binName, "--version")
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
