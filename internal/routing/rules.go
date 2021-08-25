package routing

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	errIPRuleAdd = errors.New("cannot add IP rule")
	errRulesList = errors.New("cannot list rules")
)

func (r *Routing) addIPRule(src net.IP, table, priority int) error {
	r.logger.Debug("ip rule add from " + src.String() +
		" lookup " + strconv.Itoa(table) +
		" pref " + strconv.Itoa(priority))

	rule := netlink.NewRule()
	rule.Src = netlink.NewIPNet(src)
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

func (r *Routing) deleteIPRule(src net.IP, table, priority int) error {
	r.logger.Debug("ip rule del from " + src.String() +
		" lookup " + strconv.Itoa(table) +
		" pref " + strconv.Itoa(priority))

	rule := netlink.NewRule()
	rule.Src = netlink.NewIPNet(src)
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
