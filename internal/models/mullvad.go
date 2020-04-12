package models

import "net"

type MullvadServer struct {
	IPs         []net.IP
	Country     MullvadCountry
	City        MullvadCity
	Provider    MullvadProvider
	Owned       bool
	DefaultPort uint16
}
