package settings

import "github.com/qdm12/private-internet-access-docker/internal/constants"

type DNS struct {
	Enabled           bool
	Provider          constants.DNSProvider
	AllowedHostnames  []string
	BlockMalicious    bool
	BlockSurveillance bool
	BlockAds          bool
}
