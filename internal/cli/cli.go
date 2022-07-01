package cli

type CLI struct {
	repoServersPath string
}

func New() *CLI {
	return &CLI{
		repoServersPath: "./internal/storage/servers.json",
	}
}
