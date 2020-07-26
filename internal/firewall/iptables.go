package firewall

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
)

func appendOrDelete(remove bool) string {
	if remove {
		return "--delete"
	}
	return "--append"
}

// flipRule changes an append rule in a delete rule or a delete rule into an
// append rule.
func flipRule(rule string) string {
	switch {
	case strings.HasPrefix(rule, "-A"):
		return strings.Replace(rule, "-A", "-D", 1)
	case strings.HasPrefix(rule, "--append"):
		return strings.Replace(rule, "--append", "-D", 1)
	case strings.HasPrefix(rule, "-D"):
		return strings.Replace(rule, "-D", "-A", 1)
	case strings.HasPrefix(rule, "--delete"):
		return strings.Replace(rule, "--delete", "-A", 1)
	}
	return rule
}

// Version obtains the version of the installed iptables
func (c *configurator) Version(ctx context.Context) (string, error) {
	output, err := c.commander.Run(ctx, "iptables", "--version")
	if err != nil {
		return "", err
	}
	words := strings.Fields(output)
	if len(words) < 2 {
		return "", fmt.Errorf("iptables --version: output is too short: %q", output)
	}
	return words[1], nil
}

func (c *configurator) runIptablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runIptablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *configurator) runIptablesInstruction(ctx context.Context, instruction string) error {
	c.iptablesMutex.Lock() // only one iptables command at once
	defer c.iptablesMutex.Unlock()
	if c.debug {
		fmt.Printf("iptables %s\n", instruction)
	}
	flags := strings.Fields(instruction)
	if output, err := c.commander.Run(ctx, "iptables", flags...); err != nil {
		return fmt.Errorf("failed executing \"iptables %s\": %s: %w", instruction, output, err)
	}
	return nil
}

func (c *configurator) clearAllRules(ctx context.Context) error {
	return c.runIptablesInstructions(ctx, []string{
		"--flush",        // flush all chains
		"--delete-chain", // delete all chains
	})
}

func (c *configurator) setAllPolicies(ctx context.Context, policy string) error {
	switch policy {
	case "ACCEPT", "DROP":
	default:
		return fmt.Errorf("policy %q not recognized", policy)
	}
	return c.runIptablesInstructions(ctx, []string{
		fmt.Sprintf("--policy INPUT %s", policy),
		fmt.Sprintf("--policy OUTPUT %s", policy),
		fmt.Sprintf("--policy FORWARD %s", policy),
	})
}

func (c *configurator) acceptInputThroughInterface(ctx context.Context, intf string, remove bool) error {
	return c.runIptablesInstruction(ctx, fmt.Sprintf(
		"%s INPUT -i %s -j ACCEPT", appendOrDelete(remove), intf,
	))
}

func (c *configurator) acceptOutputThroughInterface(ctx context.Context, intf string, remove bool) error {
	return c.runIptablesInstruction(ctx, fmt.Sprintf(
		"%s OUTPUT -o %s -j ACCEPT", appendOrDelete(remove), intf,
	))
}

func (c *configurator) acceptEstablishedRelatedTraffic(ctx context.Context, remove bool) error {
	return c.runIptablesInstructions(ctx, []string{
		fmt.Sprintf("%s OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT", appendOrDelete(remove)),
		fmt.Sprintf("%s INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT", appendOrDelete(remove)),
	})
}

func (c *configurator) acceptOutputTrafficToVPN(ctx context.Context, defaultInterface string, connection models.OpenVPNConnection, remove bool) error {
	return c.runIptablesInstruction(ctx,
		fmt.Sprintf("%s OUTPUT -d %s -o %s -p %s -m %s --dport %d -j ACCEPT",
			appendOrDelete(remove), connection.IP, defaultInterface, connection.Protocol, connection.Protocol, connection.Port))
}

func (c *configurator) acceptInputFromSubnetToSubnet(ctx context.Context, intf string, sourceSubnet, destinationSubnet net.IPNet, remove bool) error {
	interfaceFlag := "-i " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}
	return c.runIptablesInstruction(ctx, fmt.Sprintf(
		"%s INPUT %s -s %s -d %s -j ACCEPT", appendOrDelete(remove), interfaceFlag, sourceSubnet.String(), destinationSubnet.String(),
	))
}

// Thanks to @npawelek
func (c *configurator) acceptOutputFromSubnetToSubnet(ctx context.Context, intf string, sourceSubnet, destinationSubnet net.IPNet, remove bool) error {
	interfaceFlag := "-o " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}
	return c.runIptablesInstruction(ctx, fmt.Sprintf(
		"%s OUTPUT %s -s %s -d %s -j ACCEPT", appendOrDelete(remove), interfaceFlag, sourceSubnet.String(), destinationSubnet.String(),
	))
}

// Used for port forwarding, with intf set to tun
func (c *configurator) acceptInputToPort(ctx context.Context, intf string, port uint16, remove bool) error {
	interfaceFlag := "-i " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}
	return c.runIptablesInstructions(ctx, []string{
		fmt.Sprintf("%s INPUT %s -p tcp --dport %d -j ACCEPT", appendOrDelete(remove), interfaceFlag, port),
		fmt.Sprintf("%s INPUT %s -p udp --dport %d -j ACCEPT", appendOrDelete(remove), interfaceFlag, port),
	})
}

func (c *configurator) runUserPostRules(ctx context.Context, filepath string, remove bool) error {
	exists, err := c.fileManager.FileExists(filepath)
	if err != nil {
		return err
	} else if !exists {
		return nil
	}
	b, err := c.fileManager.ReadFile(filepath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	successfulRules := []string{}
	defer func() {
		// transaction-like rollback
		if err == nil || ctx.Err() != nil {
			return
		}
		for _, rule := range successfulRules {
			_ = c.runIptablesInstruction(ctx, flipRule(rule))
		}
	}()
	for _, line := range lines {
		if !strings.HasPrefix(line, "iptables ") {
			continue
		}
		rule := strings.TrimPrefix(line, "iptables ")
		if remove {
			rule = flipRule(rule)
		}
		if err = c.runIptablesInstruction(ctx, rule); err != nil {
			return fmt.Errorf("cannot run custom rule: %w", err)
		}
		successfulRules = append(successfulRules, rule)
	}
	return nil
}
