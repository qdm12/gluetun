package wireguard

import (
	"errors"
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

var (
	errGetLink    = errors.New("cannot get link")
	errAddAddress = errors.New("cannot add address")
)

func addAddresses(iface string, addresses []*net.IPNet) (err error) {
	link, err := netlink.LinkByName(iface)
	if err != nil {
		return fmt.Errorf("%w: %s", errGetLink, err)
	}

	for _, ipNet := range addresses {
		address := &netlink.Addr{
			IPNet: ipNet,
		}

		err = netlink.AddrAdd(link, address)
		if err != nil {
			return fmt.Errorf("%w: %s", errAddAddress, err)
		}
	}

	return nil
}
