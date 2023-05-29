//go:build !linux && !darwin

package netlink

func (n *NetLink) LinkList() (links []Link, err error) {
	panic("not implemented")
}

func (n *NetLink) LinkByName(name string) (link Link, err error) {
	panic("not implemented")
}

func (n *NetLink) LinkByIndex(index int) (link Link, err error) {
	panic("not implemented")
}

func (n *NetLink) LinkAdd(link Link) (linkIndex int, err error) {
	panic("not implemented")
}

func (n *NetLink) LinkDel(link Link) (err error) {
	panic("not implemented")
}

func (n *NetLink) LinkSetUp(link Link) (linkIndex int, err error) {
	panic("not implemented")
}

func (n *NetLink) LinkSetDown(link Link) (err error) {
	panic("not implemented")
}
