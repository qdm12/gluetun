package firewall

import (
	"context"
	"errors"
	"fmt"
)

var errRuleNotFound = errors.New("rule not found")

func (c *Config) replaceIptablesRule(ctx context.Context, oldInstruction, newInstruction string) error {
	targetRule, err := parseIptablesInstruction(oldInstruction)
	if err != nil {
		return fmt.Errorf("parsing iptables command to replace: %w", err)
	}

	lineNumber, err := findLineNumber(ctx, c.ipTables, targetRule, c.runner, c.logger)
	if err != nil {
		return fmt.Errorf("finding to-be-replaced chain rule line number: %w", err)
	} else if lineNumber == 0 {
		return fmt.Errorf("%w: matching to-be-replaced instruction %q", errRuleNotFound, oldInstruction)
	}
	parsed, err := parseIptablesInstruction(newInstruction)
	if err != nil {
		return fmt.Errorf("parsing replacement iptables command: %w", err)
	}
	parsed.operation = opReplace
	parsed.lineNumber = lineNumber
	return c.runIptablesInstruction(ctx, parsed.String())
}

func (c *Config) replaceIP6tablesRule(ctx context.Context, oldInstruction, newInstruction string) error {
	targetRule, err := parseIptablesInstruction(oldInstruction)
	if err != nil {
		return fmt.Errorf("parsing iptables command to replace: %w", err)
	}

	lineNumber, err := findLineNumber(ctx, c.ip6Tables, targetRule, c.runner, c.logger)
	if err != nil {
		return fmt.Errorf("finding to-be-replaced chain rule line number: %w", err)
	} else if lineNumber == 0 {
		return fmt.Errorf("%w: matching to-be-replaced instruction %q", errRuleNotFound, oldInstruction)
	}
	parsed, err := parseIptablesInstruction(newInstruction)
	if err != nil {
		return fmt.Errorf("parsing replacement iptables command: %w", err)
	}
	parsed.operation = opReplace
	parsed.lineNumber = lineNumber
	return c.runIP6tablesInstruction(ctx, parsed.String())
}
