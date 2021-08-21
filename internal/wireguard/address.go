package wireguard

import (
	"fmt"
	"net"

	"github.com/vishvananda/netlink"
)

func (w *Wireguard) addAddresses(link netlink.Link,
	addresses []*net.IPNet) (err error) {
	for _, ipNet := range addresses {
		address := &netlink.Addr{
			IPNet: ipNet,
		}

		err = w.netlink.AddrAdd(link, address)
		if err != nil {
			return fmt.Errorf("%w: when adding address %s to link %s",
				err, address, link.Attrs().Name)
		}
	}

	return nil
}
