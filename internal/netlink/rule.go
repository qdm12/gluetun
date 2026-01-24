package netlink

import (
	"fmt"
	"net/netip"

	"github.com/jsimonetti/rtnetlink"
)

type Rule struct {
	Priority *uint32
	Family   uint8
	Table    uint32
	Mark     *uint32
	Src      netip.Prefix
	Dst      netip.Prefix
	Flags    uint32
	Action   uint8
}

func (r *Rule) fromMessage(message rtnetlink.RuleMessage) {
	table := uint32(message.Table)
	if table == 0 || table == rtTableCompat {
		table = *message.Attributes.Table
	}
	r.Priority = message.Attributes.Priority
	r.Family = message.Family
	r.Table = table
	r.Mark = message.Attributes.FwMark
	r.Src = ipAndLengthToPrefix(message.Attributes.Src, message.SrcLength)
	r.Dst = ipAndLengthToPrefix(message.Attributes.Dst, message.DstLength)
	r.Flags = message.Flags
	r.Action = message.Action
}

func (r Rule) message() *rtnetlink.RuleMessage {
	src, srcLength := prefixToIPAndLength(r.Src)
	dst, dstLength := prefixToIPAndLength(r.Dst)

	message := &rtnetlink.RuleMessage{
		Family:    r.Family,
		SrcLength: srcLength,
		DstLength: dstLength,
		Flags:     r.Flags,
		Action:    r.Action,
		Attributes: &rtnetlink.RuleAttributes{
			Priority: r.Priority,
			FwMark:   r.Mark,
			Src:      src,
			Dst:      dst,
		},
	}

	if r.Table <= uint32(^uint8(0)) {
		message.Table = uint8(r.Table)
	} else {
		message.Table = rtTableCompat
		message.Attributes.Table = &r.Table
	}

	return message
}

func (r Rule) String() string {
	from := "all"
	if r.Src.IsValid() && !r.Src.Addr().IsUnspecified() {
		from = r.Src.String()
	}

	to := "all"
	if r.Dst.IsValid() && !r.Dst.Addr().IsUnspecified() {
		to = r.Dst.String()
	}

	priority := ""
	if r.Priority != nil {
		priority = fmt.Sprintf(" %d", *r.Priority)
	}

	return fmt.Sprintf("ip rule%s: from %s to %s table %d",
		priority, from, to, r.Table)
}

func (r Rule) debugMessage(add bool) (debugMessage string) {
	debugMessage = "ip"

	switch r.Family {
	case FamilyV4:
		debugMessage += " -f inet"
	case FamilyV6:
		debugMessage += " -f inet6"
	default:
		debugMessage += " -f " + fmt.Sprint(r.Family)
	}

	debugMessage += " rule"

	if add {
		debugMessage += " add"
	} else {
		debugMessage += " del"
	}

	if r.Src.IsValid() {
		debugMessage += " from " + r.Src.String()
	}

	if r.Dst.IsValid() {
		debugMessage += " to " + r.Dst.String()
	}

	if r.Table != 0 {
		debugMessage += " lookup " + fmt.Sprint(r.Table)
	}

	if r.Priority != nil {
		debugMessage += " pref " + fmt.Sprint(*r.Priority)
	}

	return debugMessage
}
