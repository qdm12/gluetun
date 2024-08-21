package openvpn

import "os/exec"

type CmdStarter interface {
	Start(cmd *exec.Cmd) (
		stdoutLines, stderrLines <-chan string,
		waitError <-chan error, startErr error)
}

type CmdRunStarter interface {
	Run(cmd *exec.Cmd) (output string, err error)
	CmdStarter
}
