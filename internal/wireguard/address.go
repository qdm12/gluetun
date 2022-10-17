package wireguard

import (
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

func (w *Wireguard) addAddresses(link netlink.Link,
	addresses []*net.IPNet) (err error) {
	for _, ipNet := range addresses {
		ipNetIsIPv6 := ipNet.IP.To4() == nil
		if !*w.settings.IPv6 && ipNetIsIPv6 {
			continue
		}

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
