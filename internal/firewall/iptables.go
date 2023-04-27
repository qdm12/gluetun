package firewall

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
	"os"
	"os/exec"
	"strings"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/command"
)

var (
	ErrIPTablesVersionTooShort = errors.New("iptables version string is too short")
	ErrPolicyUnknown           = errors.New("unknown policy")
	ErrNeedIP6Tables           = errors.New("ip6tables is required, please upgrade your kernel to support it")
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

// Version obtains the version of the installed iptables.
func Version(ctx context.Context, runner command.Runner) (string, error) {
	cmd := exec.CommandContext(ctx, "iptables", "--version")
	output, err := runner.Run(cmd)
	if err != nil {
		return "", err
	}
	words := strings.Fields(output)
	const minWords = 2
	if len(words) < minWords {
		return "", fmt.Errorf("%w: %s", ErrIPTablesVersionTooShort, output)
	}
	return words[1], nil
}

func (c *Config) runIptablesInstructions(ctx context.Context, instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runIptablesInstruction(ctx, instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) runIptablesInstruction(ctx context.Context, instruction string) error {
	c.iptablesMutex.Lock() // only one iptables command at once
	defer c.iptablesMutex.Unlock()

	c.logger.Debug(c.ipTables + " " + instruction)

	flags := strings.Fields(instruction)
	cmd := exec.CommandContext(ctx, c.ipTables, flags...) // #nosec G204
	if output, err := c.runner.Run(cmd); err != nil {
		return fmt.Errorf("command failed: \"%s %s\": %s: %w",
			c.ipTables, instruction, output, err)
	}
	return nil
}

func (c *Config) clearAllRules(ctx context.Context) error {
	return c.runMixedIptablesInstructions(ctx, []string{
		"--flush",        // flush all chains
		"--delete-chain", // delete all chains
	})
}

func (c *Config) setIPv4AllPolicies(ctx context.Context, policy string) error {
	switch policy {
	case "ACCEPT", "DROP":
	default:
		return fmt.Errorf("%w: %s", ErrPolicyUnknown, policy)
	}
	return c.runIptablesInstructions(ctx, []string{
		"--policy INPUT " + policy,
		"--policy OUTPUT " + policy,
		"--policy FORWARD " + policy,
	})
}

func (c *Config) acceptInputThroughInterface(ctx context.Context, intf string, remove bool) error {
	return c.runMixedIptablesInstruction(ctx, fmt.Sprintf(
		"%s INPUT -i %s -j ACCEPT", appendOrDelete(remove), intf,
	))
}

func (c *Config) acceptInputToSubnet(ctx context.Context, intf string,
	destination netip.Prefix, remove bool) error {
	interfaceFlag := "-i " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}

	instruction := fmt.Sprintf("%s INPUT %s -d %s -j ACCEPT",
		appendOrDelete(remove), interfaceFlag, destination.String())

	if destination.Addr().Is4() {
		return c.runIptablesInstruction(ctx, instruction)
	}
	if c.ip6Tables == "" {
		return fmt.Errorf("accept input to subnet %s: %w", destination, ErrNeedIP6Tables)
	}
	return c.runIP6tablesInstruction(ctx, instruction)
}

func (c *Config) acceptOutputThroughInterface(ctx context.Context, intf string, remove bool) error {
	return c.runMixedIptablesInstruction(ctx, fmt.Sprintf(
		"%s OUTPUT -o %s -j ACCEPT", appendOrDelete(remove), intf,
	))
}

func (c *Config) acceptEstablishedRelatedTraffic(ctx context.Context, remove bool) error {
	return c.runMixedIptablesInstructions(ctx, []string{
		fmt.Sprintf("%s OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT", appendOrDelete(remove)),
		fmt.Sprintf("%s INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT", appendOrDelete(remove)),
	})
}

func (c *Config) acceptOutputTrafficToVPN(ctx context.Context,
	defaultInterface string, connection models.Connection, remove bool) error {
	instruction := fmt.Sprintf("%s OUTPUT -d %s -o %s -p %s -m %s --dport %d -j ACCEPT",
		appendOrDelete(remove), connection.IP, defaultInterface, connection.Protocol,
		connection.Protocol, connection.Port)
	isIPv4 := connection.IP.To4() != nil
	if isIPv4 {
		return c.runIptablesInstruction(ctx, instruction)
	} else if c.ip6Tables == "" {
		return fmt.Errorf("accept output to VPN server: %w", ErrNeedIP6Tables)
	}
	return c.runIP6tablesInstruction(ctx, instruction)
}

// Thanks to @npawelek.
func (c *Config) acceptOutputFromIPToSubnet(ctx context.Context,
	intf string, sourceIP net.IP, destinationSubnet netip.Prefix, remove bool) error {
	doIPv4 := sourceIP.To4() != nil && destinationSubnet.Addr().Is4()

	interfaceFlag := "-o " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}

	instruction := fmt.Sprintf("%s OUTPUT %s -s %s -d %s -j ACCEPT",
		appendOrDelete(remove), interfaceFlag, sourceIP.String(), destinationSubnet.String())

	if doIPv4 {
		return c.runIptablesInstruction(ctx, instruction)
	} else if c.ip6Tables == "" {
		return fmt.Errorf("accept output from %s to %s: %w", sourceIP, destinationSubnet, ErrNeedIP6Tables)
	}
	return c.runIP6tablesInstruction(ctx, instruction)
}

// NDP uses multicast address (theres no broadcast in IPv6 like ARP uses in IPv4).
func (c *Config) acceptIpv6MulticastOutput(ctx context.Context,
	intf string, remove bool) error {
	interfaceFlag := "-o " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}
	instruction := fmt.Sprintf("%s OUTPUT %s -d ff02::1:ff/104 -j ACCEPT",
		appendOrDelete(remove), interfaceFlag)
	return c.runIP6tablesInstruction(ctx, instruction)
}

// Used for port forwarding, with intf set to tun.
func (c *Config) acceptInputToPort(ctx context.Context, intf string, port uint16, remove bool) error {
	interfaceFlag := "-i " + intf
	if intf == "*" { // all interfaces
		interfaceFlag = ""
	}
	return c.runMixedIptablesInstructions(ctx, []string{
		fmt.Sprintf("%s INPUT %s -p tcp --dport %d -j ACCEPT", appendOrDelete(remove), interfaceFlag, port),
		fmt.Sprintf("%s INPUT %s -p udp --dport %d -j ACCEPT", appendOrDelete(remove), interfaceFlag, port),
	})
}

func (c *Config) runUserPostRules(ctx context.Context, filepath string, remove bool) error {
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}
	b, err := io.ReadAll(file)
	if err != nil {
		_ = file.Close()
		return err
	}
	if err := file.Close(); err != nil {
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
		var ipv4 bool
		var rule string
		switch {
		case strings.HasPrefix(line, "iptables "):
			ipv4 = true
			rule = strings.TrimPrefix(line, "iptables ")
		case strings.HasPrefix(line, "iptables-nft "):
			ipv4 = true
			rule = strings.TrimPrefix(line, "iptables-nft ")
		case strings.HasPrefix(line, "ip6tables "):
			ipv4 = false
			rule = strings.TrimPrefix(line, "ip6tables ")
		case strings.HasPrefix(line, "ip6tables-nft "):
			ipv4 = false
			rule = strings.TrimPrefix(line, "ip6tables-nft ")
		default:
			continue
		}

		if remove {
			rule = flipRule(rule)
		}

		switch {
		case ipv4:
			err = c.runIptablesInstruction(ctx, rule)
		case c.ip6Tables == "":
			err = fmt.Errorf("running user ip6tables rule: %w", ErrNeedIP6Tables)
		default: // ipv6
			err = c.runIP6tablesInstruction(ctx, rule)
		}
		if err != nil {
			return err
		}

		successfulRules = append(successfulRules, rule)
	}
	return nil
}
