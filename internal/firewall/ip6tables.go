package firewall

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/qdm12/golibs/command"
)

var (
	ErrIP6NotSupported = errors.New("ip6tables not supported")
)

func ip6tablesSupported(ctx context.Context, runner command.Runner) (supported bool) {
	cmd := exec.CommandContext(ctx, "ip6tables", "-L")
	if _, err := runner.Run(cmd); err != nil {
		return false
	}
	return true
}

func (c *Config) runIP6tablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runIP6tablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) runIP6tablesInstruction(ctx context.Context, instruction string) error {
	if !c.ip6Tables {
		return nil
	}
	c.ip6tablesMutex.Lock() // only one ip6tables command at once
	defer c.ip6tablesMutex.Unlock()

	c.logger.Debug("ip6tables " + instruction)

	flags := strings.Fields(instruction)
	cmd := exec.CommandContext(ctx, "ip6tables", flags...)
	if output, err := c.runner.Run(cmd); err != nil {
		return fmt.Errorf("command failed: \"ip6tables %s\": %s: %w", instruction, output, err)
	}
	return nil
}

var ErrPolicyNotValid = errors.New("policy is not valid")

func (c *Config) setIPv6AllPolicies(ctx context.Context, policy string) error {
	switch policy {
	case "ACCEPT", "DROP":
	default:
		return fmt.Errorf("%w: %s", ErrPolicyNotValid, policy)
	}
	return c.runIP6tablesInstructions(ctx, []string{
		"--policy INPUT " + policy,
		"--policy OUTPUT " + policy,
		"--policy FORWARD " + policy,
	})
}
