package settings

import "github.com/qdm12/private-internet-access-docker/internal/constants"

// DNS contains settings to configure Unbound for DNS over TLS operation
type DNS struct {
	Enabled           bool
	Provider          constants.DNSProvider
	AllowedHostnames  []string
	BlockMalicious    bool
	BlockSurveillance bool
	BlockAds          bool
}
