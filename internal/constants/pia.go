package constants

import "fmt"

// PIAEncryption defines the level of encryption for communication with PIA servers
type PIAEncryption string

const (
	// PIAEncryptionNormal is the normal level of encryption for communication with PIA servers
	PIAEncryptionNormal PIAEncryption = "normal"
	// PIAEncryptionStrong is the strong level of encryption for communication with PIA servers
	PIAEncryptionStrong = "strong"
)

// ParsePIAEncryption parses a string to obtain a PIAEncryption
func ParsePIAEncryption(s string) (PIAEncryption, error) {
	switch s {
	case "normal":
		return PIAEncryptionNormal, nil
	case "strong":
		return PIAEncryptionStrong, nil
	default:
		return "", fmt.Errorf("%q can only be \"normal\" or \"strong\"", s)
	}
}

// PIARegion contains the list of regions available for PIA
type PIARegion string

// TODO add regions
