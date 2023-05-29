//go:build !linux

package netlink

func NewRule() Rule {
	return Rule{}
}

func (n *NetLink) RuleList(family int) (rules []Rule, err error) {
	panic("not implemented")
}

func (n *NetLink) RuleAdd(rule Rule) error {
	panic("not implemented")
}

func (n *NetLink) RuleDel(rule Rule) error {
	panic("not implemented")
}
