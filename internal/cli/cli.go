// Package cli defines an interface CLI to run command line operations.
package cli

var _ CLIer = (*CLI)(nil)

type CLIer interface {
	ClientKeyFormatter
	HealthChecker
	OpenvpnConfigMaker
	Updater
}

type CLI struct {
	repoServersPath string
}

func New() *CLI {
	return &CLI{
		repoServersPath: "./internal/constants/servers.json",
	}
}
