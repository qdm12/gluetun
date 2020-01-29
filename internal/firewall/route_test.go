package firewall

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/qdm12/golibs/files/mocks"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func Test_getDefaultInterface(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data             []byte
		readErr          error
		defaultInterface string
		gateway          string
		netMask          string
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
			gateway:          "172.17.0.1",
			netMask:          "16"},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			fileManager := &mocks.FileManager{}
			fileManager.On("ReadFile", constants.NetRoute).
				Return(tc.data, tc.readErr).Once()
			defaultInterface, gateway, netMask, err := getDefaultInterface(fileManager)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.defaultInterface, defaultInterface)
			assert.Equal(t, tc.gateway, gateway)
			assert.Equal(t, tc.netMask, netMask)
		})
	}
}

func Test_reversedHexToIP(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		reversedHex string
		IP          string
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
			IP:          "172.17.0.1",
			err:         nil},
		"correct hex 2": {
			reversedHex: "000011AC",
			IP:          "172.17.0.0",
			err:         nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			IP, err := reversedHexToIP(tc.reversedHex)
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
		decMask   string
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
			decMask:   "16",
			err:       nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			decMask, err := hexMaskToDecMask(tc.hexString)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.decMask, decMask)
		})
	}
}
