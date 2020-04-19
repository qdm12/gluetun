package firewall

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/models"
)

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
	flags := strings.Fields(instruction)
	if output, err := c.commander.Run(ctx, "iptables", flags...); err != nil {
		return fmt.Errorf("failed executing %q: %s: %w", instruction, output, err)
	}
	return nil
}

func (c *configurator) Clear(ctx context.Context) error {
	c.logger.Info("clearing all rules")
	return c.runIptablesInstructions(ctx, []string{
		"--flush",
		"--delete-chain",
		"-t nat --flush",
		"-t nat --delete-chain",
	})
}

func (c *configurator) AcceptAll(ctx context.Context) error {
	c.logger.Info("accepting all traffic")
	return c.runIptablesInstructions(ctx, []string{
		"-P INPUT ACCEPT",
		"-P OUTPUT ACCEPT",
		"-P FORWARD ACCEPT",
	})
}

func (c *configurator) BlockAll(ctx context.Context) error {
	c.logger.Info("blocking all traffic")
	return c.runIptablesInstructions(ctx, []string{
		"-P INPUT DROP",
		"-F OUTPUT",
		"-P OUTPUT DROP",
		"-P FORWARD DROP",
	})
}

func (c *configurator) CreateGeneralRules(ctx context.Context) error {
	c.logger.Info("creating general rules")
	return c.runIptablesInstructions(ctx, []string{
		"-A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT",
		"-A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT",
		"-A OUTPUT -o lo -j ACCEPT",
		"-A INPUT -i lo -j ACCEPT",
	})
}

func (c *configurator) CreateVPNRules(ctx context.Context, dev models.VPNDevice, defaultInterface string, connections []models.OpenVPNConnection) error {
	for _, connection := range connections {
		c.logger.Info("allowing output traffic to VPN server %s through %s on port %s %d",
			connection.IP, defaultInterface, connection.Protocol, connection.Port)
		if err := c.runIptablesInstruction(ctx,
			fmt.Sprintf("-A OUTPUT -d %s -o %s -p %s -m %s --dport %d -j ACCEPT",
				connection.IP, defaultInterface, connection.Protocol, connection.Protocol, connection.Port)); err != nil {
			return err
		}
	}
	if err := c.runIptablesInstruction(ctx, fmt.Sprintf("-A OUTPUT -o %s -j ACCEPT", dev)); err != nil {
		return err
	}
	return nil
}

func (c *configurator) CreateLocalSubnetsRules(ctx context.Context, subnet net.IPNet, extraSubnets []net.IPNet, defaultInterface string) error {
	subnetStr := subnet.String()
	c.logger.Info("accepting input and output traffic for %s", subnetStr)
	if err := c.runIptablesInstructions(ctx, []string{
		fmt.Sprintf("-A INPUT -s %s -d %s -j ACCEPT", subnetStr, subnetStr),
		fmt.Sprintf("-A OUTPUT -s %s -d %s -j ACCEPT", subnetStr, subnetStr),
	}); err != nil {
		return err
	}
	for _, extraSubnet := range extraSubnets {
		extraSubnetStr := extraSubnet.String()
		c.logger.Info("accepting input traffic through %s from %s to %s", defaultInterface, extraSubnetStr, subnetStr)
		if err := c.runIptablesInstruction(ctx,
			fmt.Sprintf("-A INPUT -i %s -s %s -d %s -j ACCEPT", defaultInterface, extraSubnetStr, subnetStr)); err != nil {
			return err
		}
		// Thanks to @npawelek
		c.logger.Info("accepting output traffic through %s from %s to %s", defaultInterface, subnetStr, extraSubnetStr)
		if err := c.runIptablesInstruction(ctx,
			fmt.Sprintf("-A OUTPUT -o %s -s %s -d %s -j ACCEPT", defaultInterface, subnetStr, extraSubnetStr)); err != nil {
			return err
		}
	}
	return nil
}

// Used for port forwarding
func (c *configurator) AllowInputTrafficOnPort(ctx context.Context, device models.VPNDevice, port uint16) error {
	c.logger.Info("accepting input traffic through %s on port %d", device, port)
	return c.runIptablesInstructions(ctx, []string{
		fmt.Sprintf("-A INPUT -i %s -p tcp --dport %d -j ACCEPT", device, port),
		fmt.Sprintf("-A INPUT -i %s -p udp --dport %d -j ACCEPT", device, port),
	})
}

func (c *configurator) AllowAnyIncomingOnPort(ctx context.Context, port uint16) error {
	c.logger.Info("accepting any input traffic on port %d", port)
	return c.runIptablesInstructions(ctx, []string{
		fmt.Sprintf("-A INPUT -p tcp --dport %d -j ACCEPT", port),
		fmt.Sprintf("-A INPUT -p udp --dport %d -j ACCEPT", port),
	})
}
