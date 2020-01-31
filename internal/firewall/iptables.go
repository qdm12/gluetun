package firewall

import (
	"fmt"
	"net"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

func (c *configurator) runIptablesInstructions(instructions []string) error {
	for _, instruction := range instructions {
		if err := c.runIptablesInstruction(instruction); err != nil {
			return err
		}
	}
	return nil
}

func (c *configurator) runIptablesInstruction(instruction string) error {
	flags := strings.Fields(instruction)
	if _, err := c.command.Run("iptables", flags...); err != nil {
		return fmt.Errorf("failed executing %q: %w", instruction, err)
	}
	return nil
}

func (c *configurator) Clear() error {
	c.logger.Info("clearing iptables rules")
	return c.runIptablesInstructions([]string{
		"--flush",
		"--delete-chain",
		"-t nat --flush",
		"-t nat --delete-chain",
	})
}

func (c *configurator) BlockAll() error {
	c.logger.Info("blocking all traffic")
	return c.runIptablesInstructions([]string{
		"-P INPUT DROP",
		"-F OUTPUT",
		"-P OUTPUT DROP",
		"-P FORWARD DROP",
	})
}

func (c *configurator) CreateGeneralRules() error {
	c.logger.Info("creating general rules")
	return c.runIptablesInstructions([]string{
		"-A OUTPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT",
		"-A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT",
		"-A OUTPUT -o lo -j ACCEPT",
		"-A INPUT -i lo -j ACCEPT",
	})
}

func (c *configurator) CreateVPNRules(dev models.VPNDevice, serverIPs []net.IP,
	defaultInterface string, port uint16, protocol constants.NetworkProtocol) error {
	for _, serverIP := range serverIPs {
		c.logger.Info("allowing output traffic to VPN server %q through %q on port %s %d",
			serverIP, defaultInterface, protocol, port)
		if err := c.runIptablesInstruction(
			fmt.Sprintf("-A OUTPUT -d %s -o %s -p %s -m %s --dport %d -j ACCEPT",
				serverIP, defaultInterface, protocol, protocol, port)); err != nil {
			return err
		}
	}
	if err := c.runIptablesInstruction(fmt.Sprintf("-A OUTPUT -o %s -j ACCEPT", dev)); err != nil {
		return err
	}
	return nil
}

func (c *configurator) CreateLocalSubnetsRules(subnet net.IPNet, extraSubnets []net.IPNet, defaultInterface string) error {
	c.logger.Info("accepting input and output traffic for %s", subnet)
	if err := c.runIptablesInstructions([]string{
		fmt.Sprintf("-A INPUT -s %s -d %s -j ACCEPT", subnet, subnet),
		fmt.Sprintf("-A OUTPUT -s %s -d %s -j ACCEPT", subnet, subnet),
	}); err != nil {
		return err
	}
	for _, extraSubnet := range extraSubnets {
		c.logger.Info("accepting input traffic through %s from %s to %s", defaultInterface, extraSubnet, subnet)
		if err := c.runIptablesInstruction(
			fmt.Sprintf("-A INPUT -i %s -s %s -d %s -j ACCEPT", defaultInterface, extraSubnet, subnet)); err != nil {
			return err
		}
	}
	return nil
}
