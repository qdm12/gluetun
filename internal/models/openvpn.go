package models

import "net"

type OpenVPNConnection struct {
	IP       net.IP
	Port     uint16
	Protocol NetworkProtocol
}

func (o *OpenVPNConnection) Equal(other OpenVPNConnection) bool {
	return o.IP.Equal(other.IP) && o.Port == other.Port && o.Protocol == other.Protocol
}
