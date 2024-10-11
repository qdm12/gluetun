//go:build !linux && !darwin

package netlink

func (n *NetLink) AddrList(link Link, family int) (
	addresses []Addr, err error,
) {
	panic("not implemented")
}

func (n *NetLink) AddrReplace(Link, Addr) error {
	panic("not implemented")
}
