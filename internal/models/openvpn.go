package models

import (
	"net"
	"strconv"
)

type OpenVPNConnection struct {
	IP       net.IP `json:"ip"`
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
	// Hostname is used for IPVanish, IVPN, Privado
	// and Windscribe for TLS verification
	Hostname string `json:"hostname"`
}

func (o *OpenVPNConnection) Equal(other OpenVPNConnection) bool {
	return o.IP.Equal(other.IP) && o.Port == other.Port && o.Protocol == other.Protocol &&
		o.Hostname == other.Hostname
}

func (o OpenVPNConnection) RemoteLine() (line string) {
	return "remote " + o.IP.String() + " " + strconv.Itoa(int(o.Port))
}

func (o OpenVPNConnection) ProtoLine() (line string) {
	return "proto " + o.Protocol
}
