package settings

import "net"

// Firewall contains settings to customize the firewall operation
type Firewall struct {
	AllowedSubnets []*net.IPNet
}
