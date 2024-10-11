package openvpn

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/qdm12/gluetun/internal/constants/openvpn"
)

var ErrVersionUnknown = errors.New("OpenVPN version is unknown")

const (
	binOpenvpn25 = "openvpn2.5"
	binOpenvpn26 = "openvpn2.6"
)

func start(ctx context.Context, starter CmdStarter, version string, flags []string) (
	stdoutLines, stderrLines <-chan string, waitError <-chan error, err error,
) {
	var bin string
	switch version {
	case openvpn.Openvpn25:
		bin = binOpenvpn25
	case openvpn.Openvpn26:
		bin = binOpenvpn26
	default:
		return nil, nil, nil, fmt.Errorf("%w: %s", ErrVersionUnknown, version)
	}

	args := []string{"--config", configPath}
	args = append(args, flags...)
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return starter.Start(cmd)
}
