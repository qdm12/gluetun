package cli

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrCommandAlreadyRegistered = errors.New("command already registered")
	errCommandUnknown           = errors.New("command is unknown")
)

type CLI struct {
	subCommands map[string]SubCommand
}

func New(subCommands ...SubCommand) (*CLI, error) {
	cli := &CLI{
		subCommands: make(map[string]SubCommand),
	}
	for _, command := range subCommands {
		if _, exists := cli.subCommands[command.Name()]; exists {
			return nil, fmt.Errorf("%w: %s", ErrCommandAlreadyRegistered, command.Name())
		}
		cli.subCommands[command.Name()] = command
	}
	return cli, nil
}

func (c *CLI) RunCommand(ctx context.Context, command string) error {
	if subCommand, exists := c.subCommands[command]; exists {
		return subCommand.Run(ctx)
	}
	c.help()
	if command == "help" {
		return nil
	}
	return fmt.Errorf("%w: %s", errCommandUnknown, command)
}

func (c *CLI) help() {
	//nolint:lll
	fmt.Printf("Usage: gluetun [COMMAND [OPTIONS]]\n\nLightweight swiss-army-knife-like VPN client to multiple VPN service providers.\n\nCommands:\n")
	for _, command := range c.subCommands {
		fmt.Printf("  %-20s\t%s\n", command.Name(), command.Description())
	}
}
