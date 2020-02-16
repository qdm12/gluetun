package models

import "net"

type MullvadServer struct {
	Country     MullvadCountry
	City        MullvadCity
	Provider    MullvadProvider
	Owned       bool
	IPs         []net.IP
	DefaultPort uint16
}
