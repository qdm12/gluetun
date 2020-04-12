package routing

import (
	"encoding/hex"
	"fmt"
	"net"
	"strconv"

	"strings"
)

type routingEntry struct {
	iface       string
	destination net.IP
	gateway     net.IP
	flags       string
	refCount    int
	use         int
	metric      int
	mask        net.IPMask
	mtu         int
	window      int
	irtt        int
}

func parseRoutingEntry(s string) (r routingEntry, err error) {
	wrapError := func(err error) error {
		return fmt.Errorf("line %q: %w", s, err)
	}
	fields := strings.Fields(s)
	if len(fields) < 11 {
		return r, wrapError(fmt.Errorf("not enough fields"))
	}
	r.iface = fields[0]
	r.destination, err = reversedHexToIPv4(fields[1])
	if err != nil {
		return r, wrapError(err)
	}
	r.gateway, err = reversedHexToIPv4(fields[2])
	if err != nil {
		return r, wrapError(err)
	}
	r.flags = fields[3]
	r.refCount, err = strconv.Atoi(fields[4])
	if err != nil {
		return r, wrapError(err)
	}
	r.use, err = strconv.Atoi(fields[5])
	if err != nil {
		return r, wrapError(err)
	}
	r.metric, err = strconv.Atoi(fields[6])
	if err != nil {
		return r, wrapError(err)
	}
	r.mask, err = hexToIPv4Mask(fields[7])
	if err != nil {
		return r, wrapError(err)
	}
	r.mtu, err = strconv.Atoi(fields[8])
	if err != nil {
		return r, wrapError(err)
	}
	r.window, err = strconv.Atoi(fields[9])
	if err != nil {
		return r, wrapError(err)
	}
	r.irtt, err = strconv.Atoi(fields[10])
	if err != nil {
		return r, wrapError(err)
	}
	return r, nil
}

func reversedHexToIPv4(reversedHex string) (ip net.IP, err error) {
	bytes, err := hex.DecodeString(reversedHex)
	if err != nil {
		return nil, fmt.Errorf("cannot parse reversed IP hex %q: %s", reversedHex, err)
	} else if len(bytes) != 4 {
		return nil, fmt.Errorf("hex string contains %d bytes instead of 4", len(bytes))
	}
	return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
}

func hexToIPv4Mask(hexString string) (mask net.IPMask, err error) {
	bytes, err := hex.DecodeString(hexString)
	if err != nil {
		return nil, fmt.Errorf("cannot parse hex mask %q: %s", hexString, err)
	} else if len(bytes) != 4 {
		return nil, fmt.Errorf("hex string contains %d bytes instead of 4", len(bytes))
	}
	return []byte{bytes[3], bytes[2], bytes[1], bytes[0]}, nil
}
