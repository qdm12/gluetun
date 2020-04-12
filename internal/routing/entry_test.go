package routing

import (
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseRoutingEntry(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s   string
		r   routingEntry
		err error
	}{
		"empty string": {
			err: fmt.Errorf("line \"\": not enough fields"),
		},
		"not enough fields": {
			s:   "a b c d e",
			err: fmt.Errorf("line \"a b c d e\": not enough fields"),
		},
		"bad destination": {
			s:   "eth0   x  0100000A  0003   0 0 0  00FFFFFF   0 0  0",
			err: fmt.Errorf("line \"eth0   x  0100000A  0003   0 0 0  00FFFFFF   0 0  0\": cannot parse reversed IP hex \"x\": encoding/hex: invalid byte: U+0078 'x'"),
		},
		"bad gateway": {
			s:   "eth0   0002A8C0  x  0003   0 0 0  00FFFFFF   0 0  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  x  0003   0 0 0  00FFFFFF   0 0  0\": cannot parse reversed IP hex \"x\": encoding/hex: invalid byte: U+0078 'x'"),
		},
		"bad ref count": {
			s:   "eth0   0002A8C0  0100000A  0003   x 0 0  00FFFFFF   0 0  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   x 0 0  00FFFFFF   0 0  0\": strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		"bad use": {
			s:   "eth0   0002A8C0  0100000A  0003   0 x 0  00FFFFFF   0 0  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   0 x 0  00FFFFFF   0 0  0\": strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		"bad metric": {
			s:   "eth0   0002A8C0  0100000A  0003   0 0 x  00FFFFFF   0 0  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   0 0 x  00FFFFFF   0 0  0\": strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		"bad mask": {
			s:   "eth0   0002A8C0  0100000A  0003   0 0 0  x   0 0  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   0 0 0  x   0 0  0\": cannot parse hex mask \"x\": encoding/hex: invalid byte: U+0078 'x'"),
		},
		"bad mtu": {
			s:   "eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   x 0  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   x 0  0\": strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		"bad window": {
			s:   "eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 x  0",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 x  0\": strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		"bad irtt": {
			s:   "eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  x",
			err: fmt.Errorf("line \"eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  x\": strconv.Atoi: parsing \"x\": invalid syntax"),
		},
		"success": {
			s: "eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  0",
			r: routingEntry{
				iface:       "eth0",
				destination: net.IP{192, 168, 2, 0},
				gateway:     net.IP{10, 0, 0, 1},
				flags:       "0003",
				mask:        net.IPMask{255, 255, 255, 0},
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r, err := parseRoutingEntry(tc.s)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.r, r)
			}
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
