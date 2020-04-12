package constants

import (
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	// PrivateInternetAccess is a VPN provider
	PrivateInternetAccess models.VPNProvider = "private internet access"
	// Mullvad is a VPN provider
	Mullvad models.VPNProvider = "mullvad"
	// Windscribe is a VPN provider
	Windscribe models.VPNProvider = "windscribe"
)

const (
	// TCP is a network protocol (reliable and slower than UDP)
	TCP models.NetworkProtocol = "tcp"
	// UDP is a network protocol (unreliable and faster than TCP)
	UDP models.NetworkProtocol = "udp"
)
