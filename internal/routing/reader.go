package routing

import (
	"bytes"
	"net"

	"fmt"
	"strings"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func parseRoutingTable(data []byte) (entries []routingEntry, err error) {
	lines := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
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

func (r *routing) routeExists(subnet net.IPNet) (exists bool, err error) {
	data, err := r.fileManager.ReadFile(string(constants.NetRoute))
	if err != nil {
		return false, fmt.Errorf("cannot check route existence: %w", err)
	}
	entries, err := parseRoutingTable(data)
	if err != nil {
		return false, fmt.Errorf("cannot check route existence: %w", err)
	}
	for _, entry := range entries {
		entrySubnet := net.IPNet{IP: entry.destination, Mask: entry.mask}
		if entrySubnet.String() == subnet.String() {
			return true, nil
		}
	}
	return false, nil
}

func (r *routing) CurrentPublicIP(defaultInterface string) (ip net.IP, err error) {
	data, err := r.fileManager.ReadFile(string(constants.NetRoute))
	if err != nil {
		return nil, fmt.Errorf("cannot find current IP address: %w", err)
	}
	entries, err := parseRoutingTable(data)
	if err != nil {
		return nil, fmt.Errorf("cannot find current IP address: %w", err)
	}
	for _, entry := range entries {
		if entry.iface == defaultInterface &&
			!ipIsPrivate(entry.destination) &&
			bytes.Equal(entry.mask, net.IPMask{255, 255, 255, 255}) {
			return entry.destination, nil
		}
	}
	return nil, fmt.Errorf("cannot find current IP address from ip routes")
}

func ipIsPrivate(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}
	privateCIDRBlocks := [8]string{
		"127.0.0.0/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}
	for i := range privateCIDRBlocks {
		_, CIDR, _ := net.ParseCIDR(privateCIDRBlocks[i])
		if CIDR.Contains(ip) {
			return true
		}
	}
	return false
}
