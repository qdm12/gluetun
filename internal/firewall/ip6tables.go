package firewall

import (
	"context"
	"fmt"
	"strings"
)

func (c *configurator) runIP6tablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runIP6tablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *configurator) runIP6tablesInstruction(ctx context.Context, instruction string) error {
	c.ip6tablesMutex.Lock() // only one ip6tables command at once
	defer c.ip6tablesMutex.Unlock()
	if c.debug {
		fmt.Println("ip6tables " + instruction)
	}
	flags := strings.Fields(instruction)
	if output, err := c.commander.Run(ctx, "ip6tables", flags...); err != nil {
		return fmt.Errorf("failed executing \"ip6tables %s\": %s: %w", instruction, output, err)
	}
	return nil
}

func (c *configurator) setIPv6AllPolicies(ctx context.Context, policy string) error {
	switch policy {
	case "ACCEPT", "DROP":
	default:
		return fmt.Errorf("policy %q not recognized", policy)
	}
	return c.runIP6tablesInstructions(ctx, []string{
		"--policy INPUT " + policy,
		"--policy OUTPUT " + policy,
		"--policy FORWARD " + policy,
	})
}
