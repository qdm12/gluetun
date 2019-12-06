package constants

type VPNProvider uint8

const (
	PrivateInternetAccess VPNProvider = iota
	Mullvad
	Windscribe
)

type Protocol uint8

const (
	TCP Protocol = iota
	UDP
)

type Region string

const ()
