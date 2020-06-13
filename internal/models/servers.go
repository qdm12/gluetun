package models

import "net"

type PIAServer struct {
	IPs    []net.IP
	Region string
}

type MullvadServer struct {
	IPs         []net.IP
	Country     string
	City        string
	ISP         string
	Owned       bool
	DefaultPort uint16
}

type WindscribeServer struct {
	Region string
	IPs    []net.IP
}

type SurfsharkServer struct {
	Region string
	IPs    []net.IP
}

type CyberghostServer struct {
	Region string
	Group  string
	IPs    []net.IP
}
