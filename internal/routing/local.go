package routing

import (
	"errors"
	"fmt"
	"net"
	"net/netip"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	ErrLinkLocalNotFound     = errors.New("local link not found")
	ErrSubnetDefaultNotFound = errors.New("default subnet not found")
	ErrSubnetLocalNotFound   = errors.New("local subnet not found")
)

type LocalNetwork struct {
	IPNet         netip.Prefix
	InterfaceName string
	IP            net.IP
}

func (r *Routing) LocalNetworks() (localNetworks []LocalNetwork, err error) {
	links, err := r.netLinker.LinkList()
	if err != nil {
		return localNetworks, fmt.Errorf("listing links: %w", err)
	}

	localLinks := make(map[int]struct{})

	for _, link := range links {
		if link.Attrs().EncapType != "ether" {
			continue
		}

		localLinks[link.Attrs().Index] = struct{}{}
		r.logger.Info("local ethernet link found: " + link.Attrs().Name)
	}

	if len(localLinks) == 0 {
		return localNetworks, fmt.Errorf("%w: in %d links", ErrLinkLocalNotFound, len(links))
	}

	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return localNetworks, fmt.Errorf("listing routes: %w", err)
	}

	for _, route := range routes {
		if route.Gw != nil || route.Dst == nil {
			continue
		} else if _, ok := localLinks[route.LinkIndex]; !ok {
			continue
		}

		var localNet LocalNetwork

		localNet.IPNet = netIPNetToNetipPrefix(*route.Dst)
		r.logger.Info("local ipnet found: " + localNet.IPNet.String())

		link, err := r.netLinker.LinkByIndex(route.LinkIndex)
		if err != nil {
			return localNetworks, fmt.Errorf("finding link at index %d: %w", route.LinkIndex, err)
		}

		localNet.InterfaceName = link.Attrs().Name

		family := netlink.FAMILY_V6
		if localNet.IPNet.Addr().Is4() {
			family = netlink.FAMILY_V4
		}
		ip, err := r.assignedIP(localNet.InterfaceName, family)
		if err != nil {
			return localNetworks, err
		}

		localNet.IP = ip

		localNetworks = append(localNetworks, localNet)
	}

	if len(localNetworks) == 0 {
		return localNetworks, fmt.Errorf("%w: in %d routes", ErrSubnetLocalNotFound, len(routes))
	}

	return localNetworks, nil
}

func (r *Routing) AddLocalRules(subnets []LocalNetwork) (err error) {
	for _, subnet := range subnets {
		// The main table is a built-in value for Linux, see "man 8 ip-route"
		const mainTable = 254

		// Local has higher priority then outbound(99) and inbound(100) as the
		// local routes might be necessary to reach the outbound/inbound routes.
		const localPriority = 98

		// Main table was setup correctly by Docker, just need to add rules to use it
		err = r.addIPRule(nil, &subnet.IPNet, mainTable, localPriority)
		if err != nil {
			return fmt.Errorf("adding rule: %v: %w", subnet.IPNet, err)
		}
	}
	return nil
}
