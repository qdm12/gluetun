package settings

import "github.com/qdm12/private-internet-access-docker/internal/constants"

type PIA struct {
	User           string
	Password       string
	Encryption     constants.PIAEncryption
	Protocol       constants.Protocol
	Region         constants.Region
	PortForwarding PortForwarding
}

type PortForwarding struct {
	Enabled  bool
	Filepath string
}
