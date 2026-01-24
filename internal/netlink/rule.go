package netlink

import (
	"fmt"

	"github.com/vishvananda/netlink"
)

func NewRule() Rule {
	// defaults found from netlink.NewRule() for fields we use,
	// the rest of the defaults is set when converting from a `Rule`
	// to a `netlink.Rule`
	return Rule{
		Priority: -1,
		Mark:     0,
	}
}

func ruleToNetlinkRule(rule Rule) (netlinkRule netlink.Rule) {
	netlinkRule = *netlink.NewRule()
	netlinkRule.Priority = rule.Priority
	netlinkRule.Family = rule.Family
	netlinkRule.Table = rule.Table
	netlinkRule.Mark = rule.Mark
	netlinkRule.Src = netipPrefixToIPNet(rule.Src)
	netlinkRule.Dst = netipPrefixToIPNet(rule.Dst)
	netlinkRule.Invert = rule.Invert
	return netlinkRule
}

func netlinkRuleToRule(netlinkRule netlink.Rule) (rule Rule) {
	return Rule{
		Priority: netlinkRule.Priority,
		Family:   netlinkRule.Family,
		Table:    netlinkRule.Table,
		Mark:     netlinkRule.Mark,
		Src:      netIPNetToNetipPrefix(netlinkRule.Src),
		Dst:      netIPNetToNetipPrefix(netlinkRule.Dst),
		Invert:   netlinkRule.Invert,
	}
}

func ruleDbgMsg(add bool, rule Rule) (debugMessage string) {
	debugMessage = "ip"

	switch rule.Family {
	case FamilyV4:
		debugMessage += " -f inet"
	case FamilyV6:
		debugMessage += " -f inet6"
	default:
		debugMessage += " -f " + fmt.Sprint(rule.Family)
	}

	debugMessage += " rule"

	if add {
		debugMessage += " add"
	} else {
		debugMessage += " del"
	}

	if rule.Src.IsValid() {
		debugMessage += " from " + rule.Src.String()
	}

	if rule.Dst.IsValid() {
		debugMessage += " to " + rule.Dst.String()
	}

	if rule.Table != 0 {
		debugMessage += " lookup " + fmt.Sprint(rule.Table)
	}

	if rule.Priority != -1 {
		debugMessage += " pref " + fmt.Sprint(rule.Priority)
	}

	return debugMessage
}
