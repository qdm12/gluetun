package routing

import (
	"bytes"
	"fmt"
	"net"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (r *Routing) addIPRule(src, dst *netip.Prefix, table, priority int) error {
	const add = true
	r.logger.Debug(ruleDbgMsg(add, src, dst, table, priority))

	rule := netlink.NewRule()
	rule.Src = NetipPrefixToIPNet(src)
	rule.Dst = NetipPrefixToIPNet(dst)
	rule.Priority = priority
	rule.Table = table

	existingRules, err := r.netLinker.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("listing rules: %w", err)
	}
	for i := range existingRules {
		if !rulesAreEqual(&existingRules[i], rule) {
			continue
		}
		return nil // already exists
	}

	if err := r.netLinker.RuleAdd(rule); err != nil {
		return fmt.Errorf("adding rule %s: %w", rule, err)
	}
	return nil
}

func (r *Routing) deleteIPRule(src, dst *netip.Prefix, table, priority int) error {
	const add = false
	r.logger.Debug(ruleDbgMsg(add, src, dst, table, priority))

	rule := netlink.NewRule()
	rule.Src = NetipPrefixToIPNet(src)
	rule.Dst = NetipPrefixToIPNet(dst)
	rule.Priority = priority
	rule.Table = table

	existingRules, err := r.netLinker.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("listing rules: %w", err)
	}
	for i := range existingRules {
		if !rulesAreEqual(&existingRules[i], rule) {
			continue
		}
		if err := r.netLinker.RuleDel(rule); err != nil {
			return fmt.Errorf("deleting rule %s: %w", rule, err)
		}
	}
	return nil
}

func ruleDbgMsg(add bool, src, dst *netip.Prefix,
	table, priority int) (debugMessage string) {
	debugMessage = "ip rule"

	if add {
		debugMessage += " add"
	} else {
		debugMessage += " del"
	}

	if src != nil {
		debugMessage += " from " + src.String()
	}

	if dst != nil {
		debugMessage += " to " + dst.String()
	}

	if table != 0 {
		debugMessage += " lookup " + fmt.Sprint(table)
	}

	if priority != -1 {
		debugMessage += " pref " + fmt.Sprint(priority)
	}

	return debugMessage
}

func rulesAreEqual(a, b *netlink.Rule) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return ipNetsAreEqual(a.Src, b.Src) &&
		ipNetsAreEqual(a.Dst, b.Dst) &&
		a.Priority == b.Priority &&
		a.Table == b.Table
}

func ipNetsAreEqual(a, b *net.IPNet) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.IP.Equal(b.IP) && bytes.Equal(a.Mask, b.Mask)
}
