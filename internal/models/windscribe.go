package models

import "net"

type WindscribeServer struct {
	Region WindscribeRegion
	IPs    []net.IP
}
