package netlink

import (
	"fmt"

	"github.com/jsimonetti/rtnetlink"
	"golang.org/x/sys/unix"
)

const (
	FlagInvert    = unix.FIB_RULE_INVERT
	ActionToTable = unix.FR_ACT_TO_TBL
)

func (n *NetLink) RuleList(family uint8) (rules []Rule, err error) {
	switch family {
	case FamilyAll:
		n.debugLogger.Debug("ip -4 rule list")
		n.debugLogger.Debug("ip -6 rule list")
	case FamilyV4:
		n.debugLogger.Debug("ip -4 rule list")
	case FamilyV6:
		n.debugLogger.Debug("ip -6 rule list")
	}

	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return nil, fmt.Errorf("dialing netlink: %w", err)
	}
	defer conn.Close()

	ruleMessages, err := conn.Rule.List()
	if err != nil {
		return nil, err
	}

	rules = make([]Rule, 0, len(ruleMessages))
	for _, message := range ruleMessages {
		if family != FamilyAll && family != message.Family {
			continue
		}
		var rule Rule
		rule.fromMessage(message)
		rules = append(rules, rule)
	}
	return rules, nil
}

func (n *NetLink) RuleAdd(rule Rule) error {
	n.debugLogger.Debug(rule.debugMessage(true))

	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return fmt.Errorf("dialing netlink: %w", err)
	}
	defer conn.Close()
	return conn.Rule.Add(rule.message())
}

func (n *NetLink) RuleDel(rule Rule) error {
	n.debugLogger.Debug(rule.debugMessage(false))

	conn, err := rtnetlink.Dial(nil)
	if err != nil {
		return fmt.Errorf("dialing netlink: %w", err)
	}
	defer conn.Close()
	return conn.Rule.Delete(rule.message())
}
