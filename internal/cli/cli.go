// Package cli defines an interface CLI to run command line operations.
package cli

type CLI struct {
	repoServersPath string
}

func New() *CLI {
	return &CLI{
		repoServersPath: "./internal/storage/servers.json",
	}
}
