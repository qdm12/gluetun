// Package cli defines an interface CLI to run command line operations.
package cli

var _ CLIer = (*CLI)(nil)

type CLIer interface {
	ClientKeyFormatter
	HealthChecker
	OpenvpnConfigMaker
	Updater
	ServersFormatter
}

type CLI struct {
	repoServersPath string
}

func New() *CLI {
	return &CLI{
		repoServersPath: "./internal/storage/servers.json",
	}
}
