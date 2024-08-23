package netlink

import (
	"fmt"
	"net/netip"
)

type Addr struct {
	Network netip.Prefix
}

func (a Addr) String() string {
	return a.Network.String()
}

type Link struct {
	Type      string
	Name      string
	Index     int
	EncapType string
	MTU       uint16
}

type Route struct {
	LinkIndex int
	Dst       netip.Prefix
	Src       netip.Addr
	Gw        netip.Addr
	Priority  int
	Family    int
	Table     int
	Type      int
}

type Rule struct {
	Priority int
	Family   int
	Table    int
	Mark     uint32
	Src      netip.Prefix
	Dst      netip.Prefix
	Invert   bool
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

	return fmt.Sprintf("ip rule %d: from %s to %s table %d",
		r.Priority, from, to, r.Table)
}
