//go:build linux || darwin

package netlink

import "github.com/vishvananda/netlink"

func (n *NetLink) LinkList() (links []Link, err error) {
	netlinkLinks, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	links = make([]Link, len(netlinkLinks))
	for i := range netlinkLinks {
		links[i] = netlinkLinkToLink(netlinkLinks[i])
	}

	return links, nil
}

func (n *NetLink) LinkByName(name string) (link Link, err error) {
	netlinkLink, err := netlink.LinkByName(name)
	if err != nil {
		return Link{}, err
	}

	return netlinkLinkToLink(netlinkLink), nil
}

func (n *NetLink) LinkByIndex(index int) (link Link, err error) {
	netlinkLink, err := netlink.LinkByIndex(index)
	if err != nil {
		return Link{}, err
	}

	return netlinkLinkToLink(netlinkLink), nil
}

func (n *NetLink) LinkAdd(link Link) (linkIndex int, err error) {
	netlinkLink := linkToNetlinkLink(&link)
	err = netlink.LinkAdd(netlinkLink)
	if err != nil {
		return 0, err
	}
	return netlinkLink.Attrs().Index, nil
}

func (n *NetLink) LinkDel(link Link) (err error) {
	return netlink.LinkDel(linkToNetlinkLink(&link))
}

func (n *NetLink) LinkSetUp(link Link) (linkIndex int, err error) {
	netlinkLink := linkToNetlinkLink(&link)
	err = netlink.LinkSetUp(netlinkLink)
	if err != nil {
		return 0, err
	}
	return netlinkLink.Attrs().Index, nil
}

func (n *NetLink) LinkSetDown(link Link) (err error) {
	return netlink.LinkSetDown(linkToNetlinkLink(&link))
}

type netlinkLinkImpl struct {
	attrs    *netlink.LinkAttrs
	linkType string
}

func (n *netlinkLinkImpl) Attrs() *netlink.LinkAttrs {
	return n.attrs
}

func (n *netlinkLinkImpl) Type() string {
	return n.linkType
}

func netlinkLinkToLink(netlinkLink netlink.Link) Link {
	attributes := netlinkLink.Attrs()
	return Link{
		Type:      netlinkLink.Type(),
		Name:      attributes.Name,
		Index:     attributes.Index,
		EncapType: attributes.EncapType,
		MTU:       uint16(attributes.MTU), //nolint:gosec
	}
}

// Warning: we must return `netlink.Link` and not `netlinkLinkImpl`
// so that the vishvananda/netlink package can compare the returned
// value against an untyped nil.
func linkToNetlinkLink(link *Link) netlink.Link {
	if link == nil {
		return nil
	}
	return &netlinkLinkImpl{
		linkType: link.Type,
		attrs: &netlink.LinkAttrs{
			Name:      link.Name,
			Index:     link.Index,
			EncapType: link.EncapType,
			MTU:       int(link.MTU),
		},
	}
}
