//go:build !linux

package netlink

const (
	// FamilyAll is a placeholder only and should not
	// be used.
	FamilyAll = iota
	// FamilyV4 is a placeholder only and should not
	// be used.
	FamilyV4
	// FamilyV6 is a placeholder only and should not
	// be used.
	FamilyV6
)

func (n *NetLink) RuleList(family int) (rules []Rule, err error) {
	panic("not implemented")
}

func (n *NetLink) RuleAdd(rule Rule) error {
	panic("not implemented")
}

func (n *NetLink) RuleDel(rule Rule) error {
	panic("not implemented")
}

func (n *NetLink) IsWireguardSupported() bool {
	panic("not implemented")
}
