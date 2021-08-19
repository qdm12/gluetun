package openvpn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/qdm12/gluetun/internal/constants"
)

var ErrVersionUnknown = errors.New("OpenVPN version is unknown")

const (
	binOpenvpn24 = "openvpn2.4"
	binOpenvpn25 = "openvpn"
)

func (c *Configurator) start(ctx context.Context, version string, flags []string) (
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
