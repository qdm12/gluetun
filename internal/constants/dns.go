package constants

import (
	"net"

	"github.com/qdm12/gluetun/internal/models"
)

const (
	// Cloudflare is a DNS over TLS provider
	Cloudflare models.DNSProvider = "cloudflare"
	// Google is a DNS over TLS provider
	Google models.DNSProvider = "google"
	// Quad9 is a DNS over TLS provider
	Quad9 models.DNSProvider = "quad9"
	// Quadrant is a DNS over TLS provider
	Quadrant models.DNSProvider = "quadrant"
	// CleanBrowsing is a DNS over TLS provider
	CleanBrowsing models.DNSProvider = "cleanbrowsing"
	// SecureDNS is a DNS over TLS provider
	SecureDNS models.DNSProvider = "securedns"
	// LibreDNS is a DNS over TLS provider
	LibreDNS models.DNSProvider = "libredns"
)

// DNSProviderMapping returns a constant mapping of dns provider name
// to their data such as IP addresses or TLS host name.
func DNSProviderMapping() map[models.DNSProvider]models.DNSProviderData {
	return map[models.DNSProvider]models.DNSProviderData{
		Cloudflare: {
			IPs:          []net.IP{{1, 1, 1, 1}, {1, 0, 0, 1}, {0x26, 0x6, 0x47, 0x0, 0x47, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x11}, {0x26, 0x6, 0x47, 0x0, 0x47, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x01}},
			SupportsTLS:  true,
			SupportsIPv6: true,
			Host:         models.DNSHost("cloudflare-dns.com"),
		},
		Google: {
			IPs:          []net.IP{{8, 8, 8, 8}, {8, 8, 4, 4}, {0x20, 0x1, 0x48, 0x60, 0x48, 0x60, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x88, 0x88}, {0x20, 0x1, 0x48, 0x60, 0x48, 0x60, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x88, 0x44}},
			SupportsTLS:  true,
			SupportsIPv6: true,
			Host:         models.DNSHost("dns.google"),
		},
		Quad9: {
			IPs:          []net.IP{{9, 9, 9, 9}, {149, 112, 112, 112}, {0x26, 0x20, 0x0, 0xfe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xfe}, {0x26, 0x20, 0x0, 0xfe, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9}},
			SupportsTLS:  true,
			SupportsIPv6: true,
			Host:         models.DNSHost("dns.quad9.net"),
		},
		Quadrant: {
			IPs:          []net.IP{{12, 159, 2, 159}, {0x20, 0x1, 0x18, 0x90, 0x14, 0xc, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x59}},
			SupportsTLS:  true,
			SupportsIPv6: true,
			Host:         models.DNSHost("dns-tls.qis.io"),
		},
		CleanBrowsing: {
			IPs:          []net.IP{{185, 228, 168, 9}, {185, 228, 169, 9}, {0x2a, 0xd, 0x2a, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}, {0x2a, 0xd, 0x2a, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2}},
			SupportsTLS:  true,
			SupportsIPv6: true,
			Host:         models.DNSHost("security-filter-dns.cleanbrowsing.org"),
		},
		SecureDNS: {
			IPs:          []net.IP{{146, 185, 167, 43}, {0x2a, 0x3, 0xb0, 0xc0, 0x0, 0x0, 0x10, 0x10, 0x0, 0x0, 0x0, 0x0, 0xe, 0x9a, 0x30, 0x1}},
			SupportsTLS:  true,
			SupportsIPv6: true,
			Host:         models.DNSHost("dot.securedns.eu"),
		},
		LibreDNS: {
			IPs:         []net.IP{{116, 203, 115, 192}},
			SupportsTLS: true,
			Host:        models.DNSHost("dot.libredns.gr"),
		},
	}
}

// Block lists URLs
const (
	AdsBlockListHostnamesURL          models.URL = "https://raw.githubusercontent.com/qdm12/files/master/ads-hostnames.updated"
	AdsBlockListIPsURL                models.URL = "https://raw.githubusercontent.com/qdm12/files/master/ads-ips.updated"
	MaliciousBlockListHostnamesURL    models.URL = "https://raw.githubusercontent.com/qdm12/files/master/malicious-hostnames.updated"
	MaliciousBlockListIPsURL          models.URL = "https://raw.githubusercontent.com/qdm12/files/master/malicious-ips.updated"
	SurveillanceBlockListHostnamesURL models.URL = "https://raw.githubusercontent.com/qdm12/files/master/surveillance-hostnames.updated"
	SurveillanceBlockListIPsURL       models.URL = "https://raw.githubusercontent.com/qdm12/files/master/surveillance-ips.updated"
)

// DNS certificates to fetch
// TODO obtain from source directly, see qdm12/updated)
const (
	NamedRootURL models.URL = "https://raw.githubusercontent.com/qdm12/files/master/named.root.updated"
	RootKeyURL   models.URL = "https://raw.githubusercontent.com/qdm12/files/master/root.key.updated"
)
