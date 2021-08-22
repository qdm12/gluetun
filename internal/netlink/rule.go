package netlink

import "github.com/vishvananda/netlink"

func (n *NetLink) RuleAdd(rule *netlink.Rule) error {
	return netlink.RuleAdd(rule)
}

func (n *NetLink) RuleDel(rule *netlink.Rule) error {
	return netlink.RuleDel(rule)
}
