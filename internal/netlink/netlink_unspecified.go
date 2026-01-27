//go:build !linux

package netlink

const (
	// FamilyAll is a placeholder only and should not
	// be used.
	FamilyAll uint8 = iota
	// FamilyV4 is a placeholder only and should not
	// be used.
	FamilyV4
	// FamilyV6 is a placeholder only and should not
	// be used.
	FamilyV6

	// DeviceTypeEthernet is a placeholder only and should not be used.
	DeviceTypeEthernet DeviceType = 0
	// DeviceTypeLoopback is a placeholder only and should not be used.
	DeviceTypeLoopback DeviceType = 0
	// DeviceTypeNone is a placeholder only and should not be used.
	DeviceTypeNone DeviceType = 0

	// iffUp is a placeholder only and should not be used.
	iffUp = 0

	// RouteTypeUnicast is a placeholder only and should not be used.
	RouteTypeUnicast = 0
	// ScopeUniverse is a placeholder only and should not be used.
	ScopeUniverse = 0
	// ProtoStatic is a placeholder only and should not be used.
	ProtoStatic = 0

	// FlagInvert is a placeholder only and should not be used.
	FlagInvert = 0
	// ActionToTable is a placeholder only and should not be used.
	ActionToTable = 0

	// rtTableCompat is a placeholder only and should not be used.
	rtTableCompat = 0
)

func (n *NetLink) RuleList(family uint8) (rules []Rule, err error) {
	panic("not implemented")
}

func (n *NetLink) RuleAdd(rule Rule) error {
	panic("not implemented")
}

func (n *NetLink) RuleDel(rule Rule) error {
	panic("not implemented")
}

func (n *NetLink) IsWireguardSupported() (bool, error) {
	panic("not implemented")
}
