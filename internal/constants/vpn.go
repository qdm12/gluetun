package constants

type VPNProvider uint8

const (
	PrivateInternetAccess VPNProvider = iota
	Mullvad
	Windscribe
)

type NetworkProtocol uint8

const (
	TCP NetworkProtocol = iota
	UDP
)

type Region string

const ()
