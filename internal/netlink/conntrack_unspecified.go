//go:build !linux

package netlink

func (n *NetLink) FlushConntrack() error {
	panic("not implemented")
}
