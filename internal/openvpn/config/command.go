package config

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/qdm12/gluetun/internal/constants"
)

var ErrVersionUnknown = errors.New("OpenVPN version is unknown")

const (
	binOpenvpn24 = "openvpn2.4"
	binOpenvpn25 = "openvpn"
)

type Starter interface {
	Start(ctx context.Context, version string, flags []string) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
}

func (c *Configurator) Start(ctx context.Context, version string, flags []string) (
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

	args := []string{"--config", constants.OpenVPNConf}
	args = append(args, flags...)
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return c.cmder.Start(cmd)
}

type VersionGetter interface {
	Version24(ctx context.Context) (version string, err error)
	Version25(ctx context.Context) (version string, err error)
}

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
