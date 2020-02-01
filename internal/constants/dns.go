package constants

import (
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

const (
	CloudflareAddress1    models.DNSForwardAddress = "1.1.1.1@853#cloudflare-dns.com"
	CloudflareAddress2    models.DNSForwardAddress = "1.0.0.1@853#cloudflare-dns.com"
	GoogleAddress1        models.DNSForwardAddress = "8.8.8.8@853#dns.google"
	GoogleAddress2        models.DNSForwardAddress = "8.8.4.4@853#dns.google"
	Quad9Address1         models.DNSForwardAddress = "9.9.9.9@853#dns.quad9.net"
	Quad9Address2         models.DNSForwardAddress = "149.112.112.112@853#dns.quad9.net"
	QuadrantAddress       models.DNSForwardAddress = "12.159.2.159@853#dns-tls.qis.io"
	CleanBrowsingAddress1 models.DNSForwardAddress = "185.228.168.9@853#security-filter-dns.cleanbrowsing.org"
	CleanBrowsingAddress2 models.DNSForwardAddress = "185.228.169.9@853#security-filter-dns.cleanbrowsing.org"
	SecureDNSAddress      models.DNSForwardAddress = "146.185.167.43@853#dot.securedns.eu"
	LibreDNSAddress       models.DNSForwardAddress = "116.203.115.192@853#dot.libredns.gr"
)

var DNSAddressesMapping = map[models.DNSProvider][]models.DNSForwardAddress{
	Cloudflare:    []models.DNSForwardAddress{CloudflareAddress1, CloudflareAddress2},
	Google:        []models.DNSForwardAddress{GoogleAddress1, GoogleAddress2},
	Quad9:         []models.DNSForwardAddress{Quad9Address1, Quad9Address2},
	Quadrant:      []models.DNSForwardAddress{QuadrantAddress},
	CleanBrowsing: []models.DNSForwardAddress{CleanBrowsingAddress1, CleanBrowsingAddress2},
	SecureDNS:     []models.DNSForwardAddress{SecureDNSAddress},
	LibreDNS:      []models.DNSForwardAddress{LibreDNSAddress},
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
