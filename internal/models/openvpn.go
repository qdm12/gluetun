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

// UpdateEmptyWith updates each field of the connection where the value is not set.
func (o *OpenVPNConnection) UpdateEmptyWith(connection OpenVPNConnection) {
	if o.IP == nil {
		o.IP = connection.IP
	}
	if o.Port == 0 {
		o.Port = connection.Port
	}
	if o.Protocol == "" {
		o.Protocol = connection.Protocol
	}
	if o.Hostname == "" {
		o.Hostname = connection.Hostname
	}
}
