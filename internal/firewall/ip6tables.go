package firewall

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/qdm12/golibs/command"
)

// findIP6tablesSupported checks for multiple iptables implementations
// and returns the iptables path that is supported. If none work, an
// empty string path is returned.
func findIP6tablesSupported(ctx context.Context, runner command.Runner) (
	ip6tablesPath string, err error) {
	ip6tablesPath, err = checkIptablesSupport(ctx, runner, "ip6tables", "ip6tables-nft")
	if errors.Is(err, ErrIPTablesNotSupported) {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return ip6tablesPath, nil
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
	if c.ip6Tables == "" {
		return nil
	}
	c.ip6tablesMutex.Lock() // only one ip6tables command at once
	defer c.ip6tablesMutex.Unlock()

	c.logger.Debug(c.ip6Tables + " " + instruction)

	flags := strings.Fields(instruction)
	cmd := exec.CommandContext(ctx, c.ip6Tables, flags...) // #nosec G204
	if output, err := c.runner.Run(cmd); err != nil {
		return fmt.Errorf("command failed: \"%s %s\": %s: %w",
			c.ip6Tables, instruction, output, err)
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
