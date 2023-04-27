package wireguard

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
	"github.com/qdm12/gluetun/internal/routing"
)

func (w *Wireguard) addAddresses(link netlink.Link,
	addresses []netip.Prefix) (err error) {
	for _, ipNet := range addresses {
		if !*w.settings.IPv6 && ipNet.Addr().Is6() {
			continue
		}

		ipNet := ipNet
		address := &netlink.Addr{
			IPNet: routing.NetipPrefixToIPNet(&ipNet),
		}

		err = w.netlink.AddrAdd(link, address)
		if err != nil {
			return fmt.Errorf("%w: when adding address %s to link %s",
				err, address, link.Attrs().Name)
		}
	}

	return nil
}
