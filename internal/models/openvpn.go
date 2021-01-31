package models

import (
	"net"
)

type OpenVPNConnection struct {
	IP       net.IP          `json:"ip"`
	Port     uint16          `json:"port"`
	Protocol NetworkProtocol `json:"protocol"`
	Hostname string          `json:"hostname"` // Privado for tls verification
}

func (o *OpenVPNConnection) Equal(other OpenVPNConnection) bool {
	return o.IP.Equal(other.IP) && o.Port == other.Port && o.Protocol == other.Protocol &&
		o.Hostname == other.Hostname
}
