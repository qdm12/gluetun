package wireguard

import (
	"github.com/vishvananda/netlink"
)

func addRule(rulePriority, firewallMark int) (
	cleanup func() error, err error) {
	rule := netlink.NewRule()
	rule.Invert = true
	rule.Priority = rulePriority
	rule.Mark = firewallMark
	rule.Table = firewallMark
	if err := netlink.RuleAdd(rule); err != nil {
		return nil, err
	}

	cleanup = func() error {
		return netlink.RuleDel(rule)
	}
	return cleanup, nil
}
