package routing

import (
	"bytes"
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	errRulesList = errors.New("cannot list rules")
)

func (r *Routing) addIPRule(src, dst *net.IPNet, table, priority int) error {
	const add = true
	r.logger.Debug(ruleDbgMsg(add, src, dst, table, priority))

	rule := netlink.NewRule()
	rule.Src = src
	rule.Dst = dst
	rule.Priority = priority
	rule.Table = table

	rules, err := r.netLinker.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("%w: %s", errRulesList, err)
	}
	for _, existingRule := range rules {
		if existingRule.Src != nil &&
			existingRule.Src.IP.Equal(rule.Src.IP) &&
			bytes.Equal(existingRule.Src.Mask, rule.Src.Mask) &&
			existingRule.Priority == rule.Priority &&
			existingRule.Table == rule.Table {
			return nil // already exists
		}
	}

	if err := r.netLinker.RuleAdd(rule); err != nil {
		return fmt.Errorf("%w: for rule: %s", err, rule)
	}
	return nil
}

func (r *Routing) deleteIPRule(src, dst *net.IPNet, table, priority int) error {
	const add = false
	r.logger.Debug(ruleDbgMsg(add, src, dst, table, priority))

	rule := netlink.NewRule()
	rule.Src = src
	rule.Dst = dst
	rule.Priority = priority
	rule.Table = table

	rules, err := r.netLinker.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("%w: %s", errRulesList, err)
	}
	for _, existingRule := range rules {
		if existingRule.Src != nil &&
			existingRule.Src.IP.Equal(rule.Src.IP) &&
			bytes.Equal(existingRule.Src.Mask, rule.Src.Mask) &&
			existingRule.Priority == rule.Priority &&
			existingRule.Table == rule.Table {
			if err := r.netLinker.RuleDel(rule); err != nil {
				return fmt.Errorf("%w: for rule: %s", err, rule)
			}
		}
	}
	return nil
}

func ruleDbgMsg(add bool, src, dst *net.IPNet,
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
