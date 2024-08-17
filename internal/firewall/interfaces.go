package firewall

import "github.com/qdm12/golibs/command"

type Runner interface {
	Run(cmd command.ExecCmd) (output string, err error)
}

type Logger interface {
	Debug(s string)
	Info(s string)
	Error(s string)
}
