package models

import "net"

type PIAServer struct {
	IPs    []net.IP
	Region PIARegion
}
