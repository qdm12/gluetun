package routing

import (
	"errors"
	"fmt"
	"net"

	"github.com/qdm12/gluetun/internal/netlink"
)

var (
	ErrRouteDefaultNotFound = errors.New("default route not found")
)

type DefaultRoute struct {
	NetInterface string
	Gateway      net.IP
	AssignedIP   net.IP
	Family       int
}

func (d DefaultRoute) String() string {
	return fmt.Sprintf("interface %s, gateway %s, assigned IP %s and family %s",
		d.NetInterface, d.Gateway, d.AssignedIP, netlink.FamilyToString(d.Family))
}

func (r *Routing) DefaultRoutes() (defaultRoutes []DefaultRoute, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("listing routes: %w", err)
	}

	for _, route := range routes {
		if route.Dst != nil {
			continue
		}
		defaultRoute := DefaultRoute{
			Gateway: route.Gw,
			Family:  route.Family,
		}
		linkIndex := route.LinkIndex
		link, err := r.netLinker.LinkByIndex(linkIndex)
		if err != nil {
			return nil, fmt.Errorf("obtaining link by index: for default route at index %d: %w", linkIndex, err)
		}
		attributes := link.Attrs()
		defaultRoute.NetInterface = attributes.Name
		family := netlink.FAMILY_V6
		if route.Gw.To4() != nil {
			family = netlink.FAMILY_V4
		}
		defaultRoute.AssignedIP, err = r.assignedIP(defaultRoute.NetInterface, family)
		if err != nil {
			return nil, fmt.Errorf("getting assigned IP of %s: %w", defaultRoute.NetInterface, err)
		}

		r.logger.Info("default route found: " + defaultRoute.String())
		defaultRoutes = append(defaultRoutes, defaultRoute)
	}

	if len(defaultRoutes) == 0 {
		return nil, fmt.Errorf("%w: in %d route(s)", ErrRouteDefaultNotFound, len(routes))
	}

	return defaultRoutes, nil
}

func DefaultRoutesInterfaces(defaultRoutes []DefaultRoute) (interfaces []string) {
	interfaces = make([]string, len(defaultRoutes))
	for i := range defaultRoutes {
		interfaces[i] = defaultRoutes[i].NetInterface
	}
	return interfaces
}
