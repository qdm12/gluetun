package wireguard

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/netlink"
	"golang.org/x/sys/unix"
)

func (w *Wireguard) addRule4(rulePriority, firewallMark int) (
	cleanup func() error, err error) {
	rule := netlink.NewRule()
	rule.Invert = true
	rule.Priority = rulePriority
	rule.Mark = firewallMark
	rule.Table = firewallMark
	rule.Family = unix.AF_INET
	if err := w.netlink.RuleAdd(rule); err != nil {
		return nil, fmt.Errorf("cannot add rule %s: %w", rule, err)
	}

	cleanup = func() error {
		err := w.netlink.RuleDel(rule)
		if err != nil {
			return fmt.Errorf("cannot delete rule %s: %w", rule, err)
		}
		return nil
	}
	return cleanup, nil
}

func (w *Wireguard) addRule6(rulePriority, firewallMark int) (
	cleanup func() error, err error) {
	rule := netlink.NewRule()
	rule.Invert = true
	rule.Priority = rulePriority
	rule.Mark = firewallMark
	rule.Table = firewallMark
	rule.Family = unix.AF_INET6
	if err := w.netlink.RuleAdd(rule); err != nil {
		return nil, fmt.Errorf("cannot add rule %s: %w", rule, err)
	}

	cleanup = func() error {
		err := w.netlink.RuleDel(rule)
		if err != nil {
			return fmt.Errorf("cannot delete rule %s: %w", rule, err)
		}
		return nil
	}
	return cleanup, nil
}
