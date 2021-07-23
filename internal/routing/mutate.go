package routing

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/vishvananda/netlink"
)

var (
	ErrRouteReplace = errors.New("cannot replace route")
	ErrRouteDelete  = errors.New("cannot delete route")
	ErrRuleAdd      = errors.New("cannot add routing rule")
	ErrRuleDel      = errors.New("cannot delete routing rule")
)

func (r *routing) addRouteVia(destination net.IPNet, gateway net.IP, iface string, table int) error {
	destinationStr := destination.String()
	r.logger.Info("adding route for " + destinationStr)
	r.logger.Debug("ip route replace " + destinationStr +
		" via " + gateway.String() +
		" dev " + iface +
		" table " + strconv.Itoa(table))

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("%w: interface %s: %s", ErrLinkByName, iface, err)
	}
	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := netlink.RouteReplace(&route); err != nil {
		return fmt.Errorf("%w: for subnet %s at interface %s: %s",
			ErrRouteReplace, destinationStr, iface, err)
	}
	return nil
}

func (r *routing) deleteRouteVia(destination net.IPNet, gateway net.IP, iface string, table int) (err error) {
	destinationStr := destination.String()
	r.logger.Info("deleting route for " + destinationStr)
	r.logger.Debug("ip route delete " + destinationStr +
		" via " + gateway.String() +
		" dev " + iface +
		" table " + strconv.Itoa(table))

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("%w: for interface %s: %s", ErrLinkByName, iface, err)
	}
	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := netlink.RouteDel(&route); err != nil {
		return fmt.Errorf("%w: for subnet %s at interface %s: %s",
			ErrRouteDelete, destinationStr, iface, err)
	}
	return nil
}

func (r *routing) addIPRule(src net.IP, table, priority int) error {
	r.logger.Debug("ip rule add from " + src.String() +
		" lookup " + strconv.Itoa(table) +
		" pref " + strconv.Itoa(priority))

	rule := netlink.NewRule()
	rule.Src = netlink.NewIPNet(src)
	rule.Priority = priority
	rule.Table = table

	rules, err := netlink.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRulesList, err)
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

	if err := netlink.RuleAdd(rule); err != nil {
		return fmt.Errorf("%w: for rule %q: %s", ErrRuleAdd, rule, err)
	}
	return nil
}

func (r *routing) deleteIPRule(src net.IP, table, priority int) error {
	r.logger.Debug("ip rule del from " + src.String() +
		" lookup " + strconv.Itoa(table) +
		" pref " + strconv.Itoa(priority))

	rule := netlink.NewRule()
	rule.Src = netlink.NewIPNet(src)
	rule.Priority = priority
	rule.Table = table

	rules, err := netlink.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrRulesList, err)
	}
	for _, existingRule := range rules {
		if existingRule.Src != nil &&
			existingRule.Src.IP.Equal(rule.Src.IP) &&
			bytes.Equal(existingRule.Src.Mask, rule.Src.Mask) &&
			existingRule.Priority == rule.Priority &&
			existingRule.Table == rule.Table {
			if err := netlink.RuleDel(rule); err != nil {
				return fmt.Errorf("%w: for rule %q: %s", ErrRuleDel, rule, err)
			}
		}
	}
	return nil
}
