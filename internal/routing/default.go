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

type DefaultRouteGetter interface {
	DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error)
}

func (r *Routing) DefaultRoute() (defaultInterface string, defaultGateway net.IP, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return "", nil, fmt.Errorf("cannot list routes: %w", err)
	}
	for _, route := range routes {
		if route.Dst == nil {
			defaultGateway = route.Gw
			linkIndex := route.LinkIndex
			link, err := r.netLinker.LinkByIndex(linkIndex)
			if err != nil {
				return "", nil, fmt.Errorf("cannot obtain link by index: for default route at index %d: %w", linkIndex, err)
			}
			attributes := link.Attrs()
			defaultInterface = attributes.Name
			r.logger.Info("default route found: interface " + defaultInterface +
				", gateway " + defaultGateway.String())
			return defaultInterface, defaultGateway, nil
		}
	}
	return "", nil, fmt.Errorf("%w: in %d route(s)", ErrRouteDefaultNotFound, len(routes))
}

type DefaultIPGetter interface {
	DefaultIP() (defaultIP net.IP, err error)
}

func (r *Routing) DefaultIP() (ip net.IP, err error) {
	routes, err := r.netLinker.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("cannot list routes: %w", err)
	}

	defaultLinkName := ""
	for _, route := range routes {
		if route.Dst == nil {
			linkIndex := route.LinkIndex
			link, err := r.netLinker.LinkByIndex(linkIndex)
			if err != nil {
				return nil, fmt.Errorf("cannot find link by index: for default route at index %d: %w", linkIndex, err)
			}
			defaultLinkName = link.Attrs().Name
		}
	}
	if defaultLinkName == "" {
		return nil, fmt.Errorf("%w: in %d route(s)", ErrLinkDefaultNotFound, len(routes))
	}

	return r.assignedIP(defaultLinkName)
}
