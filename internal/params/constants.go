package params

/////////
// VPN //
/////////
type Protocol uint8
const (
	TCP Protocol = iota
	UDP
)

type VPNProvider uint8
const (
	PIA VPNProvider = iota
	Mullvad
	Windscribe
)

type Region string
const ()

// PIA
type PIAEncryption uint8
const (
	PIAEncryptionNormal PIAEncryption = iota
	PIAEncryptionStrong
)

// Mullvad

// Windscribe


/////////
// DNS //
/////////
type DNSProvider uint8
const (
	Cloudflare DNSProvider = iota
	Google
	Quad9
	CleanBrowsing
)

///////////////
// TINYPROXY //
///////////////
type TinyProxyLogLevel uint8
const (
	TinyProxyInfoLevel TinyProxyLogLevel = iota
	TinyProxyWarnLevel
	TinyProxyErrorLevel
	TinyProxyCriticalLevel
)
