package wireguard

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) addRule(rulePriority int, firewallMark uint32,
	family int,
) (cleanup func() error, err error) {
	rule := netlink.NewRule()
	rule.Invert = true
	rule.Priority = rulePriority
	rule.Mark = firewallMark
	rule.Table = int(firewallMark)
	rule.Family = family
	if err := w.netlink.RuleAdd(rule); err != nil {
		return nil, fmt.Errorf("adding %s: %w", rule, err)
	}

	cleanup = func() error {
		err := w.netlink.RuleDel(rule)
		if err != nil {
			return fmt.Errorf("deleting rule %s: %w", rule, err)
		}
		return nil
	}
	return cleanup, nil
}
