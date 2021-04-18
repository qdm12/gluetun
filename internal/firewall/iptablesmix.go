package firewall

import (
	"context"
)

func (c *configurator) runMixedIptablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runMixedIptablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *configurator) runMixedIptablesInstruction(ctx context.Context, instruction string) error {
	if err := c.runIptablesInstruction(ctx, instruction); err != nil {
		return err
	}
	return c.runIP6tablesInstruction(ctx, instruction)
}
