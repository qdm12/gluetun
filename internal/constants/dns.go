package constants

// DNSProvider is a DNS over TLS server provider name
type DNSProvider string

const (
	// Cloudflare is a DNS over TLS provider
	Cloudflare DNSProvider = "cloudflare"
	// Google is a DNS over TLS provider
	Google = "google"
	// Quad9 is a DNS over TLS provider
	Quad9 = "quad9"
	// Quadrant is a DNS over TLS provider
	Quadrant = "quadrant"
	// CleanBrowsing is a DNS over TLS provider
	CleanBrowsing = "cleanbrowsing"
	// SecureDNS is a DNS over TLS provider
	SecureDNS = "securedns"
	// LibreDNS is a DNS over TLS provider
	LibreDNS = "libredns"
)

// GetForwardAddresses gets forwarded addresses corresponding to a
// DNS over TLS provider.
func (p *DNSProvider) GetForwardAddresses() []string {
	switch *p {
	case Cloudflare:
		return []string{"1.1.1.1@853#cloudflare-dns.com", "1.0.0.1@853#cloudflare-dns.com"}
	case Google:
		return []string{"8.8.8.8@853#dns.google", "8.8.4.4@853#dns.google"}
	case Quad9:
		return []string{"9.9.9.9@853#dns.quad9.net", "149.112.112.112@853#dns.quad9.net"}
	case Quadrant:
		return []string{"12.159.2.159@853#dns-tls.qis.io"}
	case CleanBrowsing:
		return []string{
			"185.228.168.9@853#security-filter-dns.cleanbrowsing.org",
			"185.228.169.9@853#security-filter-dns.cleanbrowsing.org"}
	case SecureDNS:
		return []string{"146.185.167.43@853#dot.securedns.eu"}
	case LibreDNS:
		return []string{"116.203.115.192@853#dot.libredns.gr"}
	default:
		return nil
	}
}

// Block lists URLs
const (
	AdsBlockListHostnamesURL          = "https://raw.githubusercontent.com/qdm12/files/master/ads-hostnames.updated"
	AdsBlockListIPsURL                = "https://raw.githubusercontent.com/qdm12/files/master/ads-ips.updated"
	MaliciousBlockListHostnamesURL    = "https://raw.githubusercontent.com/qdm12/files/master/malicious-hostnames.updated"
	MaliciousBlockListIPsURL          = "https://raw.githubusercontent.com/qdm12/files/master/malicious-ips.updated"
	SurveillanceBlockListHostnamesURL = "https://raw.githubusercontent.com/qdm12/files/master/surveillance-hostnames.updated"
	SurveillanceBlockListIPsURL       = "https://raw.githubusercontent.com/qdm12/files/master/surveillance-ips.updated"
)

// DNS certificates to fetch (TODO obtain from source directly, see qdm12/updated)
const (
	NamedRootURL = "https://raw.githubusercontent.com/qdm12/files/master/named.root.updated"
	RootKeyURL   = "https://raw.githubusercontent.com/qdm12/files/master/root.key.updated"
)
