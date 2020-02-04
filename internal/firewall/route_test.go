package firewall

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	filesmocks "github.com/qdm12/golibs/files/mocks"
	loggingmocks "github.com/qdm12/golibs/logging/mocks"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func Test_getDefaultRoute(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data             []byte
		readErr          error
		defaultInterface string
		defaultGateway   net.IP
		defaultSubnet    net.IPNet
		err              error
	}{
		"no data": {
			err: fmt.Errorf("not enough lines (1) found in %s", constants.NetRoute)},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error")},
		"not enough fields line 1": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    00000000
eth0    000011AC        00000000        0001    0       0       0       0000FFFF        0       0       0`),
			err: fmt.Errorf("not enough fields in \"eth0    00000000\"")},
		"not enough fields line 2": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0
eth0    000011AC        00000000        0001    0       0       0`),
			err: fmt.Errorf("not enough fields in \"eth0    000011AC        00000000        0001    0       0       0\"")},
		"bad gateway": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    00000000        x               0003    0       0       0       00000000        0       0       0
eth0    000011AC        00000000        0001    0       0       0       0000FFFF        0       0       0`),
			err: fmt.Errorf("cannot parse reversed IP hex \"x\": encoding/hex: invalid byte: U+0078 'x'")},
		"bad net number": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0
eth0    x               00000000        0001    0       0       0       0000FFFF        0       0       0`),
			err: fmt.Errorf("cannot parse reversed IP hex \"x\": encoding/hex: invalid byte: U+0078 'x'")},
		"bad net mask": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0
eth0    000011AC        00000000        0001    0       0       0       x               0       0       0`),
			err: fmt.Errorf("cannot parse hex mask \"x\": encoding/hex: invalid byte: U+0078 'x'")},
		"success": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT        
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0             
eth0    000011AC        00000000        0001    0       0       0       0000FFFF        0       0       0`),
			defaultInterface: "eth0",
			defaultGateway:   net.IP{0xac, 0x11, 0x0, 0x1},
			defaultSubnet: net.IPNet{
				IP:   net.IP{0xac, 0x11, 0x0, 0x0},
				Mask: net.IPMask{0xff, 0xff, 0x0, 0x0},
			}},
	} // TODO find full subnet 172.17.0.0/16
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fileManager := &filesmocks.FileManager{}
			fileManager.On("ReadFile", string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Once()
			logger := &loggingmocks.Logger{}
			logger.On("Info", "%s: detecting default network route", logPrefix).Once()
			if tc.err == nil {
				logger.On("Info", "%s: default route found: interface %s, gateway %s, subnet %s",
					logPrefix, tc.defaultInterface, tc.defaultGateway.String(), tc.defaultSubnet.String()).Once()
			}
			c := &configurator{logger: logger, fileManager: fileManager}
			defaultInterface, defaultGateway, defaultSubnet, err := c.GetDefaultRoute()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.defaultInterface, defaultInterface)
			assert.Equal(t, tc.defaultGateway, defaultGateway)
			assert.Equal(t, tc.defaultSubnet, defaultSubnet)
			fileManager.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}

func Test_reversedHexToIPv4(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		reversedHex string
		IP          net.IP
		err         error
	}{
		"empty hex": {
			err: fmt.Errorf("hex string contains 0 bytes instead of 4")},
		"bad hex": {
			reversedHex: "x",
			err:         fmt.Errorf("cannot parse reversed IP hex \"x\": encoding/hex: invalid byte: U+0078 'x'")},
		"3 bytes hex": {
			reversedHex: "9abcde",
			err:         fmt.Errorf("hex string contains 3 bytes instead of 4")},
		"correct hex": {
			reversedHex: "010011AC",
			IP:          []byte{0xac, 0x11, 0x0, 0x1},
			err:         nil},
		"correct hex 2": {
			reversedHex: "000011AC",
			IP:          []byte{0xac, 0x11, 0x0, 0x0},
			err:         nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			IP, err := reversedHexToIPv4(tc.reversedHex)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.IP, IP)
		})
	}
}

func Test_hexMaskToDecMask(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		hexString string
		mask      net.IPMask
		err       error
	}{
		"empty hex": {
			err: fmt.Errorf("hex string contains 0 bytes instead of 4")},
		"bad hex": {
			hexString: "x",
			err:       fmt.Errorf("cannot parse hex mask \"x\": encoding/hex: invalid byte: U+0078 'x'")},
		"3 bytes hex": {
			hexString: "9abcde",
			err:       fmt.Errorf("hex string contains 3 bytes instead of 4")},
		"16": {
			hexString: "0000FFFF",
			mask:      []byte{0xff, 0xff, 0x0, 0x0},
			err:       nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mask, err := hexToIPv4Mask(tc.hexString)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.mask, mask)
		})
	}
}
