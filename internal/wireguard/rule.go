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
		if strings.HasSuffix(err.Error(), "file exists") {
			w.logger.Info("if you are using Kubernetes, this may fix the error below: " +
				"https://github.com/qdm12/gluetun-wiki/blob/main/setup/advanced/kubernetes.md#adding-ipv6-rule--file-exists")
		}
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
