package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_DNSProvider_GetForwardAddresses(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		provider           DNSProvider
		forwardedAddresses []string
	}{
		"Cloudflare": {Cloudflare, []string{"1.1.1.1@853#cloudflare-dns.com", "1.0.0.1@853#cloudflare-dns.com"}},
		"Google":     {Google, []string{"8.8.8.8@853#dns.google", "8.8.4.4@853#dns.google"}},
		"Quad9":      {Quad9, []string{"9.9.9.9@853#dns.quad9.net", "149.112.112.112@853#dns.quad9.net"}},
		"Quadrant":   {Quadrant, []string{"12.159.2.159@853#dns-tls.qis.io"}},
		"CleanBrowsing": {CleanBrowsing, []string{"185.228.168.9@853#security-filter-dns.cleanbrowsing.org",
			"185.228.169.9@853#security-filter-dns.cleanbrowsing.org"}},
		"SecureDNS": {SecureDNS, []string{"146.185.167.43@853#dot.securedns.eu"}},
		"LibreDNS":  {LibreDNS, []string{"116.203.115.192@853#dot.libredns.gr"}},
		"Unknown":   {"unknown", nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			forwardedAddresses := tc.provider.GetForwardAddresses()
			assert.Equal(t, tc.forwardedAddresses, forwardedAddresses)
		})
	}
}
