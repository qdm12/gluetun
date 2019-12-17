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
	// CleanBrowsing is a DNS over TLS provider
	CleanBrowsing = "cleanbrowsing"
	// TODO add more
)
