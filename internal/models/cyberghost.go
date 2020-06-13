package models

import "net"

type CyberghostServer struct {
	Region CyberghostRegion
	Group  CyberghostGroup
	IPs    []net.IP
}
