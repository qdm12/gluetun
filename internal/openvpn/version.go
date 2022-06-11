package openvpn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func (c *Configurator) Version24(ctx context.Context) (version string, err error) {
	return c.version(ctx, binOpenvpn24)
}

func (c *Configurator) Version25(ctx context.Context) (version string, err error) {
	return c.version(ctx, binOpenvpn25)
}

var ErrVersionTooShort = errors.New("version output is too short")

func (c *Configurator) version(ctx context.Context, binName string) (version string, err error) {
	cmd := exec.CommandContext(ctx, binName, "--version")
	output, err := c.cmder.Run(cmd)
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
