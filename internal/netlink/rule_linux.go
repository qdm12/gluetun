package netlink

import "github.com/vishvananda/netlink"

func (n *NetLink) RuleList(family int) (rules []Rule, err error) {
	switch family {
	case FamilyAll:
		n.debugLogger.Debug("ip -4 rule list")
		n.debugLogger.Debug("ip -6 rule list")
	case FamilyV4:
		n.debugLogger.Debug("ip -4 rule list")
	case FamilyV6:
		n.debugLogger.Debug("ip -6 rule list")
	}
	netlinkRules, err := netlink.RuleList(family)
	if err != nil {
		return nil, err
	}

	rules = make([]Rule, len(netlinkRules))
	for i := range netlinkRules {
		rules[i] = netlinkRuleToRule(netlinkRules[i])
	}
	return rules, nil
}

func (n *NetLink) RuleAdd(rule Rule) error {
	n.debugLogger.Debug(ruleDbgMsg(true, rule))
	netlinkRule := ruleToNetlinkRule(rule)
	return netlink.RuleAdd(&netlinkRule)
}

func (n *NetLink) RuleDel(rule Rule) error {
	n.debugLogger.Debug(ruleDbgMsg(false, rule))
	netlinkRule := ruleToNetlinkRule(rule)
	return netlink.RuleDel(&netlinkRule)
}
