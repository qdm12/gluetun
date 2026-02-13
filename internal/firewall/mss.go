package firewall

import (
	"context"
	"fmt"

	"github.com/qdm12/gluetun/internal/pmtud/constants"
)

func (c *Config) SetVPNTCPMSS(ctx context.Context, mtu uint32) error {
	c.stateMutex.Lock()
	defer c.stateMutex.Unlock()

	if !c.enabled || c.vpnIntf == "" {
		return nil
	}

	ruleTemplate := "-t mangle %s FORWARD -o " + c.vpnIntf + " -p tcp --tcp-flags SYN,RST SYN -j TCPMSS --set-mss %d"

	onlyRemove := mtu == 0

	const mysteriousOverhead = 20 // most likely TCP options, such as the 12B of timestamps

	const ipv4Overhead = constants.IPv4HeaderLength + constants.BaseTCPHeaderLength + mysteriousOverhead
	tcpMSSIPv4 := mtu - ipv4Overhead
	var instructions []string
	if c.vpnTCPMSSIPv4 != 0 {
		instructions = append(instructions, fmt.Sprintf(ruleTemplate, "--delete", c.vpnTCPMSSIPv4))
	} else if !onlyRemove {
		instructions = append(instructions, fmt.Sprintf(ruleTemplate, "--append", tcpMSSIPv4))
	}
	err := c.runIptablesInstructions(ctx, instructions)
	if err != nil {
		return fmt.Errorf("setting TCP MSS for IPv4: %w", err)
	} else if !onlyRemove {
		c.vpnTCPMSSIPv4 = tcpMSSIPv4
	}

	const ipv6Overhead = constants.IPv6HeaderLength + constants.BaseTCPHeaderLength + mysteriousOverhead
	tcpMSSIPv6 := mtu - ipv6Overhead
	instructions = []string{}
	if c.vpnTCPMSSIPv6 != 0 {
		instructions = append(instructions, fmt.Sprintf(ruleTemplate, "--delete", c.vpnTCPMSSIPv6))
	} else if !onlyRemove {
		instructions = append(instructions, fmt.Sprintf(ruleTemplate, "--append", tcpMSSIPv6))
	}
	err = c.runIP6tablesInstructions(ctx, instructions)
	if err != nil {
		return fmt.Errorf("setting TCP MSS for IPv6: %w", err)
	} else if !onlyRemove {
		c.vpnTCPMSSIPv6 = tcpMSSIPv6
	}

	return nil
}
