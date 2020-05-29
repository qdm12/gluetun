package models

import "net"

type SurfsharkServer struct {
	Region SurfsharkRegion
	IPs    []net.IP
}
