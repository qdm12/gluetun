package routing

import (
	"net"

	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func parseRoutingTable(data []byte) (entries []routingEntry, err error) {
	lines := strings.Split(string(data), "\n")
	lines = lines[1:]
	entries = make([]routingEntry, len(lines))
	for i := range lines {
		entries[i], err = parseRoutingEntry(lines[i])
		if err != nil {
			return nil, fmt.Errorf("line %d in %s: %w", i+1, constants.NetRoute, err)
		}
	}
	return entries, nil
}

func (r *routing) DefaultRoute() (defaultInterface string, defaultGateway net.IP, defaultSubnet net.IPNet, err error) {
	r.logger.Info("detecting default network route")
	data, err := r.fileManager.ReadFile(string(constants.NetRoute))
	if err != nil {
		return "", nil, defaultSubnet, err
	}
	entries, err := parseRoutingTable(data)
	if err != nil {
		return "", nil, defaultSubnet, err
	}
	if len(entries) < 2 {
		return "", nil, defaultSubnet, fmt.Errorf("not enough entries (%d) found in %s", len(entries), constants.NetRoute)
	}
	defaultInterface = entries[0].iface
	defaultGateway = entries[0].gateway
	defaultSubnet = net.IPNet{IP: entries[1].destination, Mask: entries[1].mask}
	r.logger.Info("default route found: interface %s, gateway %s, subnet %s", defaultInterface, defaultGateway.String(), defaultSubnet.String())
	return defaultInterface, defaultGateway, defaultSubnet, nil
}
