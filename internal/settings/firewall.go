package settings

import "net"

type Firewall struct {
	AllowedSubnets []*net.IPNet
}
