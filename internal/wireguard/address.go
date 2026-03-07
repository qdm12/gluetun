package wireguard

import (
	"fmt"
	"net/netip"
)

func (w *Wireguard) addAddresses(linkIndex uint32,
	addresses []netip.Prefix,
) (err error) {
	for _, address := range addresses {
		if !*w.settings.IPv6 && address.Addr().Is6() {
			continue
		}

		err = w.netlink.AddrReplace(linkIndex, address)
		if err != nil {
			return fmt.Errorf("%w: when adding address %s to link with index %d",
				err, address, linkIndex)
		}
	}

	return nil
}
