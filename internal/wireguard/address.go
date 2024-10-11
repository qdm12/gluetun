package wireguard

import (
	"fmt"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) addAddresses(link netlink.Link,
	addresses []netip.Prefix,
) (err error) {
	for _, ipNet := range addresses {
		if !*w.settings.IPv6 && ipNet.Addr().Is6() {
			continue
		}

		address := netlink.Addr{
			Network: ipNet,
		}

		err = w.netlink.AddrReplace(link, address)
		if err != nil {
			return fmt.Errorf("%w: when adding address %s to link %s",
				err, address, link.Name)
		}
	}

	return nil
}
