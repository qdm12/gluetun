package constants

import (
	"net"

	"github.com/qdm12/private-internet-access-docker/internal/models"
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
		Cloudflare: models.DNSProviderData{
			IPs:         []net.IP{{1, 1, 1, 1}, {1, 0, 0, 1}},
			SupportsTLS: true,
			Host:        models.DNSHost("cloudflare-dns.com"),
		},
		Google: models.DNSProviderData{
			IPs:         []net.IP{{8, 8, 8, 8}, {8, 8, 4, 4}},
			SupportsTLS: true,
			Host:        models.DNSHost("dns.google"),
		},
		Quad9: models.DNSProviderData{
			IPs:         []net.IP{{9, 9, 9, 9}, {149, 112, 112, 112}},
			SupportsTLS: true,
			Host:        models.DNSHost("dns.quad9.net"),
		},
		Quadrant: models.DNSProviderData{
			IPs:         []net.IP{{12, 159, 2, 159}},
			SupportsTLS: true,
			Host:        models.DNSHost("dns-tls.qis.io"),
		},
		CleanBrowsing: models.DNSProviderData{
			IPs:         []net.IP{{185, 228, 168, 9}, {185, 228, 169, 9}},
			SupportsTLS: true,
			Host:        models.DNSHost("security-filter-dns.cleanbrowsing.org"),
		},
		SecureDNS: models.DNSProviderData{
			IPs:         []net.IP{{146, 185, 167, 43}},
			SupportsTLS: true,
			Host:        models.DNSHost("dot.securedns.eu"),
		},
		LibreDNS: models.DNSProviderData{
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
