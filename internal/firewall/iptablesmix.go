package firewall

import (
	"context"
	"fmt"
)

func (c *Config) runMixedIptablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runMixedIptablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) runMixedIptablesInstruction(ctx context.Context, instruction string) error {
	if err := c.runIptablesInstruction(ctx, instruction); err != nil {
		return err
	}
	return c.runIP6tablesInstruction(ctx, instruction)
}

func (c *Config) runIPv4AndV6IptablesInstructions(ctx context.Context,
	ipv4Instructions, ipv6Instructions []string,
) error {
	if err := c.runIptablesInstructions(ctx, ipv4Instructions); err != nil {
		return fmt.Errorf("running iptables instructions: %w", err)
	}
	if err := c.runIP6tablesInstructions(ctx, ipv6Instructions); err != nil {
		return fmt.Errorf("running ip6tables instructions: %w", err)
	}
	return nil
}
