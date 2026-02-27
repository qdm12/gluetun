package iptables

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// isDeleteMatchInstruction returns true if the iptables instruction
// is a delete instruction by rule matching. It returns false if the
// instruction is a delete instruction by line number, or not a delete
// instruction.
func isDeleteMatchInstruction(instruction string) bool {
	fields := strings.Fields(instruction)
	for i, field := range fields {
		switch {
		case field != "-D" && field != "--delete":
			continue
		case i == len(fields)-1: // malformed: missing chain name
			return false
		case i == len(fields)-2: // chain name is last field
			return true
		default:
			// chain name is fields[i+1]
			const base, bitLength = 10, 16
			_, err := strconv.ParseUint(fields[i+2], base, bitLength)
			return err != nil // not a line number
		}
	}
	return false
}

func deleteIPTablesRule(ctx context.Context, iptablesBinary, instruction string,
	runner CmdRunner, logger Logger,
) (err error) {
	targetRule, err := parseIptablesInstruction(instruction)
	if err != nil {
		return fmt.Errorf("parsing iptables command: %w", err)
	}

	lineNumber, err := findLineNumber(ctx, iptablesBinary,
		targetRule, runner, logger)
	if err != nil {
		return fmt.Errorf("finding iptables chain rule line number: %w", err)
	} else if lineNumber == 0 {
		logger.Debug("rule matching \"" + instruction + "\" not found")
		return nil
	}
	logger.Debug(fmt.Sprintf("found iptables chain rule matching %q at line number %d",
		instruction, lineNumber))

	cmd := exec.CommandContext(ctx, iptablesBinary, "-t", targetRule.table,
		"-D", targetRule.chain, fmt.Sprint(lineNumber)) // #nosec G204
	logger.Debug(cmd.String())
	output, err := runner.Run(cmd)
	if err != nil {
		err = fmt.Errorf("command failed: %q: %w", cmd, err)
		if output != "" {
			err = fmt.Errorf("%w: %s", err, output)
		}
		return err
	}

	return nil
}

// findLineNumber finds the line number of an iptables rule.
// It returns 0 if the rule is not found.
func findLineNumber(ctx context.Context, iptablesBinary string,
	instruction iptablesInstruction, runner CmdRunner, logger Logger) (
	lineNumber uint16, err error,
) {
	listFlags := []string{
		"-t", instruction.table, "-L", instruction.chain,
		"--line-numbers", "-n", "-v",
	}
	cmd := exec.CommandContext(ctx, iptablesBinary, listFlags...) // #nosec G204
	logger.Debug(cmd.String())
	output, err := runner.Run(cmd)
	if err != nil {
		err = fmt.Errorf("command failed: %q: %w", cmd, err)
		if output != "" {
			err = fmt.Errorf("%w: %s", err, output)
		}
		return 0, err
	}

	chain, err := parseChain(output)
	if err != nil {
		return 0, fmt.Errorf("parsing chain list: %w", err)
	}

	for _, rule := range chain.rules {
		if instruction.equalToRule(instruction.table, chain.name, rule) {
			return rule.lineNumber, nil
		}
	}

	return 0, nil
}
