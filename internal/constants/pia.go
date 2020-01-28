package constants

// PIAEncryption defines the level of encryption for communication with PIA servers
type PIAEncryption string

const (
	// PIAEncryptionNormal is the normal level of encryption for communication with PIA servers
	PIAEncryptionNormal PIAEncryption = "normal"
	// PIAEncryptionStrong is the strong level of encryption for communication with PIA servers
	PIAEncryptionStrong = "strong"
)

// PIARegion contains the list of regions available for PIA
type PIARegion string

// TODO add regions

const (
	PIAOpenVPNURL = "https://www.privateinternetaccess.com/openvpn"
)
