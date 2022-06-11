package netlink

import "github.com/vishvananda/netlink"

type Rule = netlink.Rule

func NewRule() *Rule {
	return netlink.NewRule()
}

func (n *NetLink) RuleList(family int) (rules []Rule, err error) {
	return netlink.RuleList(family)
}

func (n *NetLink) RuleAdd(rule *Rule) error {
	return netlink.RuleAdd(rule)
}

func (n *NetLink) RuleDel(rule *Rule) error {
	return netlink.RuleDel(rule)
}
