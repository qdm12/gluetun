package firewall

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/golibs/command"
)

var (
	ErrIP6Tables       = errors.New("failed ip6tables command")
	ErrIP6NotSupported = errors.New("ip6tables not supported")
)

func ip6tablesSupported(ctx context.Context, commander command.Commander) (supported bool) {
	if _, err := commander.Run(ctx, "ip6tables", "-L"); err != nil {
		return false
	}
	return true
}

func (c *configurator) runIP6tablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runIP6tablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *configurator) runIP6tablesInstruction(ctx context.Context, instruction string) error {
	if !c.ip6Tables {
		return nil
	}
	c.ip6tablesMutex.Lock() // only one ip6tables command at once
	defer c.ip6tablesMutex.Unlock()
	if c.debug {
		fmt.Println("ip6tables " + instruction)
	}
	flags := strings.Fields(instruction)
	if output, err := c.commander.Run(ctx, "ip6tables", flags...); err != nil {
		return fmt.Errorf("%w: \"ip6tables %s\": %s: %s", ErrIP6Tables, instruction, output, err)
	}
	return nil
}

var errPolicyNotValid = errors.New("policy is not valid")

func (c *configurator) setIPv6AllPolicies(ctx context.Context, policy string) error {
	switch policy {
	case "ACCEPT", "DROP":
	default:
		return fmt.Errorf("%w: %s", errPolicyNotValid, policy)
	}
	return c.runIP6tablesInstructions(ctx, []string{
		"--policy INPUT " + policy,
		"--policy OUTPUT " + policy,
		"--policy FORWARD " + policy,
	})
}
