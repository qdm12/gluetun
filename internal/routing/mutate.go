package routing

import (
	"bytes"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func (r *routing) addRouteVia(destination net.IPNet, gateway net.IP, iface string, table int) error {
	destinationStr := destination.String()
	r.logger.Info("adding route for %s", destinationStr)
	if r.debug {
		fmt.Printf("ip route replace %s via %s dev %s table %d\n", destinationStr, gateway, iface, table)
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("cannot add route for %s: %w", destinationStr, err)
	}
	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := netlink.RouteReplace(&route); err != nil {
		return fmt.Errorf("cannot add route for %s: %w", destinationStr, err)
	}
	return nil
}

func (r *routing) deleteRouteVia(destination net.IPNet, gateway net.IP, iface string, table int) (err error) {
	destinationStr := destination.String()
	r.logger.Info("deleting route for %s", destinationStr)
	if r.debug {
		fmt.Printf("ip route delete %s via %s dev %s table %d\n", destinationStr, gateway, iface, table)
	}

	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("cannot delete route for %s: %w", destinationStr, err)
	}
	route := netlink.Route{
		Dst:       &destination,
		Gw:        gateway,
		LinkIndex: link.Attrs().Index,
		Table:     table,
	}
	if err := netlink.RouteDel(&route); err != nil {
		return fmt.Errorf("cannot delete route for %s: %w", destinationStr, err)
	}
	return nil
}

func (r *routing) addIPRule(src net.IP, table, priority int) error {
	if r.debug {
		fmt.Printf("ip rule add from %s lookup %d pref %d\n",
			src, table, priority)
	}

	rule := netlink.NewRule()
	rule.Src = netlink.NewIPNet(src)
	rule.Priority = priority
	rule.Table = table

	rules, err := netlink.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("cannot add ip rule: %w", err)
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

	return netlink.RuleAdd(rule)
}

func (r *routing) deleteIPRule(src net.IP, table, priority int) error {
	if r.debug {
		fmt.Printf("ip rule del from %s lookup %d pref %d\n",
			src, table, priority)
	}

	rule := netlink.NewRule()
	rule.Src = netlink.NewIPNet(src)
	rule.Priority = priority
	rule.Table = table

	rules, err := netlink.RuleList(netlink.FAMILY_ALL)
	if err != nil {
		return fmt.Errorf("cannot add ip rule: %w", err)
	}
	for _, existingRule := range rules {
		if existingRule.Src != nil &&
			existingRule.Src.IP.Equal(rule.Src.IP) &&
			bytes.Equal(existingRule.Src.Mask, rule.Src.Mask) &&
			existingRule.Priority == rule.Priority &&
			existingRule.Table == rule.Table {
			return netlink.RuleDel(rule)
		}
	}
	return nil
}
