package wireguard

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) addRule(rulePriority, firewallMark int) (
	cleanup func() error, err error) {
	rule := netlink.NewRule()
	rule.Invert = true
	rule.Priority = rulePriority
	rule.Mark = firewallMark
	rule.Table = firewallMark
	if err := w.netlink.RuleAdd(rule); err != nil {
		return nil, fmt.Errorf("%w: when adding rule: %s", err, rule)
	}

	cleanup = func() error {
		err := w.netlink.RuleDel(rule)
		if err != nil {
			return fmt.Errorf("%w: when deleting rule: %s", err, rule)
		}
		return nil
	}
	return cleanup, nil
}
