package firewall

import "os/exec"

type CmdRunner interface {
	Run(cmd *exec.Cmd) (output string, err error)
}

type Logger interface {
	Debug(s string)
	Info(s string)
	Error(s string)
}
