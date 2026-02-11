package firewall

import (
	"context"
	"fmt"
	"net/netip"
	"strconv"
)

func (c *Config) SetAllowedPort(ctx context.Context, port uint16, intf string) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal state")
		existingInterfaces, ok := c.allowedInputPorts[port]
		if !ok {
			existingInterfaces = make(map[string]struct{})
		}
		existingInterfaces[intf] = struct{}{}
		c.allowedInputPorts[port] = existingInterfaces
		return nil
	}

	netInterfaces, has := c.allowedInputPorts[port]
	if !has {
		netInterfaces = make(map[string]struct{})
	} else if _, exists := netInterfaces[intf]; exists {
		return nil
	}

	c.logger.Info("setting allowed input port " + fmt.Sprint(port) + " through interface " + intf + "...")

	const remove = false
	if err := c.acceptInputToPort(ctx, intf, port, remove); err != nil {
		return fmt.Errorf("allowing input to port %d through interface %s: %w",
			port, intf, err)
	}
	netInterfaces[intf] = struct{}{}
	c.allowedInputPorts[port] = netInterfaces

	return nil
}

func (c *Config) RemoveAllowedPort(ctx context.Context, port uint16) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if port == 0 {
		return nil
	}

	if !c.enabled {
		c.logger.Info("firewall disabled, only updating allowed ports internal list")
		delete(c.allowedInputPorts, port)
		return nil
	}

	c.logger.Info("removing allowed port " + strconv.Itoa(int(port)) + "...")

	interfacesSet, ok := c.allowedInputPorts[port]
	if !ok {
		return nil
	}

	const remove = true
	for netInterface := range interfacesSet {
		err := c.acceptInputToPort(ctx, netInterface, port, remove)
		if err != nil {
			return fmt.Errorf("removing allowed port %d on interface %s: %w",
				port, netInterface, err)
		}
		delete(interfacesSet, netInterface)
	}

	// All interfaces were removed successfully, so remove the port entry.
	delete(c.allowedInputPorts, port)

	return nil
}

// RestrictOutputAddrPort allows outgoing traffic to a specific IP and port for both tcp and udp,
// while blocking other tcp or udp traffic to that port going to other IP addresses, both IPv4 and IPv6.
// If the port was previously allowed for another IP address, that previous allowance will be removed.
// Giving an invalid address will remove any existing restrictions for the port specified.
func (c *Config) RestrictOutputAddrPort(ctx context.Context, addrPort netip.AddrPort) (err error) {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()
	existingIP := c.outputAddrPort[addrPort.Port()]

	switch {
	case existingIP == addrPort.Addr():
		return nil
	case !addrPort.Addr().IsValid():
		// invalid address, remove any existing rules for the port
		return c.removeOutputAddrPortRestriction(ctx, existingIP, addrPort.Port())
	case !existingIP.IsValid():
		// no previous existing address for the port
		return c.insertOutputAddrPortRestriction(ctx, addrPort)
	default:
		// existing rule in the same IP family or different family
		return c.replaceOutputAddrPortRestriction(ctx, existingIP, addrPort)
	}
}

func (c *Config) removeOutputAddrPortRestriction(ctx context.Context, existingIP netip.Addr, port uint16) (err error) {
	commonInstructions := []string{
		fmt.Sprintf("--delete OUTPUT -p udp --dport %d -j DROP", port),
		fmt.Sprintf("--delete OUTPUT -p tcp --dport %d -j DROP", port),
	}
	ipv4Instructions := commonInstructions
	ipv6Instructions := commonInstructions

	familySpecificInstructions := []string{
		fmt.Sprintf("--delete OUTPUT -p udp --dport %d -d %s -j ACCEPT", port, existingIP),
		fmt.Sprintf("--delete OUTPUT -p tcp --dport %d -d %s -j ACCEPT", port, existingIP),
	}
	if existingIP.Is4() {
		ipv4Instructions = append(ipv4Instructions, familySpecificInstructions...)
	} else {
		ipv6Instructions = append(ipv6Instructions, familySpecificInstructions...)
	}

	err = c.runIPv4AndV6IptablesInstructions(ctx, ipv4Instructions, ipv6Instructions)
	if err != nil {
		return err
	}
	delete(c.outputAddrPort, port)
	return nil
}

func (c *Config) insertOutputAddrPortRestriction(ctx context.Context, addrPort netip.AddrPort) (err error) {
	commonInstructions := []string{
		fmt.Sprintf("--insert OUTPUT -p udp --dport %d -j DROP", addrPort.Port()),
		fmt.Sprintf("--insert OUTPUT -p tcp --dport %d -j DROP", addrPort.Port()),
	}
	ipv4Instructions := commonInstructions
	ipv6Instructions := commonInstructions

	familySpecificInstructions := []string{
		fmt.Sprintf("--insert OUTPUT -p udp --dport %d -d %s -j ACCEPT", addrPort.Port(), addrPort.Addr()),
		fmt.Sprintf("--insert OUTPUT -p tcp --dport %d -d %s -j ACCEPT", addrPort.Port(), addrPort.Addr()),
	}
	if addrPort.Addr().Is4() {
		ipv4Instructions = append(ipv4Instructions, familySpecificInstructions...)
	} else {
		ipv6Instructions = append(ipv6Instructions, familySpecificInstructions...)
	}
	err = c.runIPv4AndV6IptablesInstructions(ctx, ipv4Instructions, ipv6Instructions)
	if err != nil {
		return err
	}
	c.outputAddrPort[addrPort.Port()] = addrPort.Addr()
	return nil
}

func (c *Config) replaceOutputAddrPortRestriction(ctx context.Context,
	existingIP netip.Addr, addrPort netip.AddrPort,
) (err error) {
	for _, protocol := range [...]string{"udp", "tcp"} {
		switch {
		case existingIP.Is4() && addrPort.Addr().Is4():
			oldInstruction := fmt.Sprintf("--insert OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), existingIP)
			newInstruction := fmt.Sprintf("--insert OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), addrPort.Addr())
			err = c.replaceIptablesRule(ctx, oldInstruction, newInstruction)
			if err != nil {
				return fmt.Errorf("replacing existing IPv4 rule: %w", err)
			}
		case existingIP.Is6() && addrPort.Addr().Is6():
			oldInstruction := fmt.Sprintf("--insert OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), existingIP)
			newInstruction := fmt.Sprintf("--insert OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), addrPort.Addr())
			err = c.replaceIP6tablesRule(ctx, oldInstruction, newInstruction)
			if err != nil {
				return fmt.Errorf("replacing existing IPv6 rule: %w", err)
			}
		case existingIP.Is4() && addrPort.Addr().Is6():
			instruction := fmt.Sprintf("--delete OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), existingIP)
			err = c.runIptablesInstruction(ctx, instruction)
			if err != nil {
				return fmt.Errorf("removing existing IPv4 rule: %w", err)
			}
			instruction = fmt.Sprintf("--insert OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), addrPort.Addr())
			err = c.runIP6tablesInstruction(ctx, instruction)
			if err != nil {
				return fmt.Errorf("inserting new IPv6 rule: %w", err)
			}
		case existingIP.Is6() && addrPort.Addr().Is4():
			instruction := fmt.Sprintf("--delete OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), existingIP)
			err = c.runIP6tablesInstruction(ctx, instruction)
			if err != nil {
				return fmt.Errorf("removing existing IPv6 rule: %w", err)
			}
			instruction = fmt.Sprintf("--insert OUTPUT -p %s --dport %d -d %s -j ACCEPT",
				protocol, addrPort.Port(), addrPort.Addr())
			err = c.runIptablesInstruction(ctx, instruction)
			if err != nil {
				return fmt.Errorf("inserting new IPv4 rule: %w", err)
			}
		}
	}
	c.outputAddrPort[addrPort.Port()] = addrPort.Addr()
	return nil
}
