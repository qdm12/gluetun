package routing

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (r *Routing) addIPRule(src, dst netip.Prefix, table, priority int) error {
	rule := netlink.NewRule()
	rule.Src = src
	rule.Dst = dst
	rule.Priority = priority
	rule.Table = table

	existingRules, err := r.netLinker.RuleList(netlink.FamilyAll)
	if err != nil {
		return fmt.Errorf("listing rules: %w", err)
	}
	for i := range existingRules {
		if !rulesAreEqual(existingRules[i], rule) {
			continue
		}
		return nil // already exists
	}

	if err := r.netLinker.RuleAdd(rule); err != nil {
		return fmt.Errorf("adding %s: %w", rule, err)
	}
	return nil
}

func (r *Routing) deleteIPRule(src, dst netip.Prefix, table, priority int) error {
	rule := netlink.NewRule()
	rule.Src = src
	rule.Dst = dst
	rule.Priority = priority
	rule.Table = table

	existingRules, err := r.netLinker.RuleList(netlink.FamilyAll)
	if err != nil {
		return fmt.Errorf("listing rules: %w", err)
	}
	for i := range existingRules {
		if !rulesAreEqual(existingRules[i], rule) {
			continue
		}
		if err := r.netLinker.RuleDel(rule); err != nil {
			return fmt.Errorf("deleting rule %s: %w", rule, err)
		}
	}
	return nil
}

func rulesAreEqual(a, b netlink.Rule) bool {
	return ipPrefixesAreEqual(a.Src, b.Src) &&
		ipPrefixesAreEqual(a.Dst, b.Dst) &&
		a.Priority == b.Priority &&
		a.Table == b.Table
}

func ipPrefixesAreEqual(a, b netip.Prefix) bool {
	if !a.IsValid() && !b.IsValid() {
		return true
	}
	if !a.IsValid() || !b.IsValid() {
		return false
	}
	return a.Bits() == b.Bits() &&
		a.Addr().Compare(b.Addr()) == 0
}
