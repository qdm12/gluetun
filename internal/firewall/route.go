package firewall

import (
	"encoding/hex"

	"github.com/qdm12/golibs/files"

	"fmt"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"strings"
)

func getDefaultInterface(fileManager files.FileManager) (defaultInterface, gateway, netMask string, err error) {
	data, err := fileManager.ReadFile(constants.NetRoute)
	if err != nil {
		return "", "", "", err
	}
	// Verify number of lines and fields
	lines := strings.Split(string(data), "\n")
	if len(lines) < 3 {
		return "", "", "", fmt.Errorf("not enough lines (%d) found in %s", len(lines), constants.NetRoute)
	}
	fieldsLine1 := strings.Fields(lines[1])
	if len(fieldsLine1) < 3 {
		return "", "", "", fmt.Errorf("not enough fields in %q", lines[1])
	}
	fieldsLine2 := strings.Fields(lines[2])
	if len(fieldsLine2) < 8 {
		return "", "", "", fmt.Errorf("not enough fields in %q", lines[2])
	}
	// get information
	defaultInterface = fieldsLine1[0]
	gateway, err = reversedHexToIP(fieldsLine1[2])
	if err != nil {
		return "", "", "", err
	}
	netMask, err = hexMaskToDecMask(fieldsLine2[7])
	if err != nil {
		return "", "", "", err
	}
	return defaultInterface, gateway, netMask, nil
}

func reversedHexToIP(reversedHex string) (IP string, err error) {
	bytes, err := hex.DecodeString(reversedHex)
	if err != nil {
		return "", fmt.Errorf("cannot parse reversed IP hex %q: %s", reversedHex, err)
	} else if len(bytes) != 4 {
		return "", fmt.Errorf("hex string contains %d bytes instead of 4", len(bytes))
	}
	return fmt.Sprintf("%v.%v.%v.%v", bytes[3], bytes[2], bytes[1], bytes[0]), nil
}

func hexMaskToDecMask(hexString string) (decMask string, err error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return "", fmt.Errorf("cannot parse hex mask %q: %s", hexString, err)
	} else if len(bytes) != 4 {
		return "", fmt.Errorf("hex string contains %d bytes instead of 4", len(bytes))
	}
	bitString := fmt.Sprintf("%08b", bytes)
	ones := strings.Count(bitString, "1")
	return fmt.Sprintf("%d", ones), nil
}
