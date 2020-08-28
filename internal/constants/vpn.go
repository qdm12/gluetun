package constants

import (
	"github.com/qdm12/gluetun/internal/models"
)

const (
	// PrivateInternetAccess is a VPN provider
	PrivateInternetAccess models.VPNProvider = "private internet access"
	// PrivateInternetAccessOld is the pre summer 2020 PIA provider
	PrivateInternetAccessOld models.VPNProvider = "private internet access old"
	// Mullvad is a VPN provider
	Mullvad models.VPNProvider = "mullvad"
	// Windscribe is a VPN provider
	Windscribe models.VPNProvider = "windscribe"
	// Surfshark is a VPN provider
	Surfshark models.VPNProvider = "surfshark"
	// Cyberghost is a VPN provider
	Cyberghost models.VPNProvider = "cyberghost"
	// Vyprvpn is a VPN provider
	Vyprvpn models.VPNProvider = "vyprvpn"
	// NordVPN is a VPN provider
	Nordvpn models.VPNProvider = "nordvpn"
	// PureVPN is a VPN provider
	Purevpn models.VPNProvider = "purevpn"
)

const (
	// TCP is a network protocol (reliable and slower than UDP)
	TCP models.NetworkProtocol = "tcp"
	// UDP is a network protocol (unreliable and faster than TCP)
	UDP models.NetworkProtocol = "udp"
)
