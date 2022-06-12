package ipinfo

import "net"

type Response struct {
	IP       net.IP `json:"ip,omitempty"`
	Region   string `json:"region,omitempty"`
	Country  string `json:"country,omitempty"`
	City     string `json:"city,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Loc      string `json:"loc,omitempty"`
	Org      string `json:"org,omitempty"`
	Postal   string `json:"postal,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

func (r Response) Copy() (copied Response) {
	copied = r
	copied.IP = make(net.IP, len(r.IP))
	copy(copied.IP, r.IP)
	return copied
}
