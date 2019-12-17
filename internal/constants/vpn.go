package constants

import "fmt"

// VPNProvider is the name of the VPN provider to be used
type VPNProvider string

const (
	// PrivateInternetAccess is a VPN provider
	PrivateInternetAccess VPNProvider = "private internet access"
	// Mullvad is a VPN provider
	Mullvad = "mullvad"
	// Windscribe is a VPN provider
	Windscribe = "windscribe"
)

// NetworkProtocol contains the network protocol to be used to communicate with the VPN servers
type NetworkProtocol string

const (
	// TCP is a network protocol (reliable and slower than UDP)
	TCP NetworkProtocol = "tcp"
	// UDP is a network protocol (unreliable and faster than TCP)
	UDP = "udp"
)

// ParseNetworkProtocol parses a string to obtain the network protocol to use
func ParseNetworkProtocol(s string) (NetworkProtocol, error) {
	switch s {
	case "tcp":
		return TCP, nil
	case "udp":
		return UDP, nil
	default:
		return "", fmt.Errorf("network protocol can only be \"tcp\" or \"udp\"")
	}
}
