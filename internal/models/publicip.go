package models

import "net"

type PublicIP struct {
	IP           net.IP `json:"public_ip,omitempty"`
	Region       string `json:"region,omitempty"`
	Country      string `json:"country,omitempty"`
	City         string `json:"city,omitempty"`
	Hostname     string `json:"hostname,omitempty"`
	Location     string `json:"location,omitempty"`
	Organization string `json:"organization,omitempty"`
	PostalCode   string `json:"postal_code,omitempty"`
	Timezone     string `json:"timezone,omitempty"`
}

func (p *PublicIP) Copy() (publicIPCopy PublicIP) {
	publicIPCopy = PublicIP{
		IP:           make(net.IP, len(p.IP)),
		Region:       p.Region,
		Country:      p.Country,
		City:         p.City,
		Hostname:     p.Hostname,
		Location:     p.Location,
		Organization: p.Organization,
		PostalCode:   p.PostalCode,
		Timezone:     p.Timezone,
	}
	copy(publicIPCopy.IP, p.IP)
	return publicIPCopy
}
