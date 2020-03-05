package models

import "net"

// DNSProviderData contains information for a DNS provider
type DNSProviderData struct {
	IPs          []net.IP
	SupportsTLS  bool
	SupportsIPv6 bool
	Host         DNSHost
}
