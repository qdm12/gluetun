package models

import "net"

type IPInfoData struct {
	IP       net.IP `json:"public_ip"`
	Region   string `json:"region,omitempty"`
	Country  string `json:"country,omitempty"`
	City     string `json:"city,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Loc      string `json:"loc,omitempty"`
	Org      string `json:"org,omitempty"`
	Postal   string `json:"postal,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

func (i IPInfoData) Copy() (copied IPInfoData) {
	copied = i
	copied.IP = make(net.IP, len(i.IP))
	copy(copied.IP, i.IP)
	return copied
}

func (i *IPInfoData) SetIP(ip net.IP) {
	i.IP = make(net.IP, len(ip))
	copy(i.IP, ip)
}
