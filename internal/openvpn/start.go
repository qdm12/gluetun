package openvpn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/command"
)

var ErrVersionUnknown = errors.New("OpenVPN version is unknown")

const (
	binOpenvpn24 = "openvpn2.4"
	binOpenvpn25 = "openvpn2.5"
)

func start(ctx context.Context, starter command.Starter, version string, flags []string) (
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

	args := []string{"--config", configPath}
	args = append(args, flags...)
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return starter.Start(cmd)
}
