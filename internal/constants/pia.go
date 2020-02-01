package constants

import (
	"github.com/qdm12/private-internet-access-docker/internal/models"
)

const (
	// PIAEncryptionNormal is the normal level of encryption for communication with PIA servers
	PIAEncryptionNormal models.PIAEncryption = "normal"
	// PIAEncryptionStrong is the strong level of encryption for communication with PIA servers
	PIAEncryptionStrong models.PIAEncryption = "strong"
)

// TODO add regions

const (
	PIAOpenVPNURL models.URL = "https://www.privateinternetaccess.com/openvpn"
)
