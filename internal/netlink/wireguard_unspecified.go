//go:build !linux

package netlink

func (n *NetLink) IsWireguardSupported() (ok bool, err error) {
	panic("not implemented")
}
