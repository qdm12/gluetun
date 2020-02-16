package models

import "net"

type OpenVPNConnection struct {
	IP       net.IP
	Port     uint16
	Protocol NetworkProtocol
}
