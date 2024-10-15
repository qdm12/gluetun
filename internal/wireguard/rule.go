package wireguard

import (
	"fmt"
	"strings"

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
		if strings.Contains(err.Error(), "file exists") {
			rules, listErr := w.netlink.RuleList(family)
			if listErr != nil {
				return nil, fmt.Errorf("listing rules for family %d due to %q: %w",
					family, err, listErr)
			}
			ruleStrings := make([]string, len(rules))
			for i := range rules {
				ruleStrings[i] = rules[i].String()
			}
			w.logger.Info("existing rules are:\n" + strings.Join(ruleStrings, "\n"))
		}
		return nil, fmt.Errorf("adding rule %s: %w", rule, err)
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
