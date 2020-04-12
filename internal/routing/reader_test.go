package routing

import (
	"fmt"
	"net"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

//go:generate mockgen -destination=mockLogger_test.go -package=routing github.com/qdm12/golibs/logging Logger
//go:generate mockgen -destination=mockFilemanager_test.go -package=routing github.com/qdm12/golibs/files FileManager

func Test_parseRoutingTable(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data    []byte
		entries []routingEntry
		err     error
	}{
		"nil data": {
			entries: []routingEntry{},
		},
		"legend only": {
			data:    []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT`),
			entries: []routingEntry{},
		},
		"legend and single line": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  0`),
			entries: []routingEntry{{
				iface:       "eth0",
				destination: net.IP{192, 168, 2, 0},
				gateway:     net.IP{10, 0, 0, 1},
				flags:       "0003",
				mask:        net.IPMask{255, 255, 255, 0},
			}},
		},
		"legend and two lines": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  0
eth0   0002A8C0  0100000A  0002   0 0 0  00FFFFFF   0 0  0`),
			entries: []routingEntry{
				{
					iface:       "eth0",
					destination: net.IP{192, 168, 2, 0},
					gateway:     net.IP{10, 0, 0, 1},
					flags:       "0003",
					mask:        net.IPMask{255, 255, 255, 0},
				},
				{
					iface:       "eth0",
					destination: net.IP{192, 168, 2, 0},
					gateway:     net.IP{10, 0, 0, 1},
					flags:       "0002",
					mask:        net.IPMask{255, 255, 255, 0},
				}},
		},
		"error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x  0100000A  0003   0 0 0  00FFFFFF   0 0  0`),
			entries: nil,
			err:     fmt.Errorf("line 1 in /proc/net/route: line \"eth0   x  0100000A  0003   0 0 0  00FFFFFF   0 0  0\": cannot parse reversed IP hex \"x\": encoding/hex: invalid byte: U+0078 'x'"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			entries, err := parseRoutingTable(tc.data)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.entries, entries)
		})
	}
}

func Test_DefaultRoute(t *testing.T) {
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
			err: fmt.Errorf("not enough entries (0) found in %s", constants.NetRoute)},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error")},
		"parse error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x`),
			err: fmt.Errorf("line 1 in /proc/net/route: line \"eth0   x\": not enough fields")},
		"single entry": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT        
eth0    00000000        050A090A        0003    0       0       0       00000080        0       0       0`),
			err: fmt.Errorf("not enough entries (1) found in %s", constants.NetRoute)},
		"success": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT        
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0
eth0    000011AC        00000000        0001    0       0       0       0000FFFF        0       0       0`),
			defaultInterface: "eth0",
			defaultGateway:   net.IP{172, 17, 0, 1},
			defaultSubnet: net.IPNet{
				IP:   net.IP{172, 17, 0, 0},
				Mask: net.IPMask{255, 255, 0, 0},
			}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockLogger := NewMockLogger(mockCtrl)
			mockFilemanager := NewMockFileManager(mockCtrl)

			mockFilemanager.EXPECT().ReadFile(string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Times(1)
			mockLogger.EXPECT().Info("detecting default network route").Times(1)
			if tc.err == nil {
				mockLogger.EXPECT().Info(
					"default route found: interface %s, gateway %s, subnet %s",
					tc.defaultInterface, tc.defaultGateway.String(), tc.defaultSubnet.String(),
				).Times(1)
			}
			r := &routing{logger: mockLogger, fileManager: mockFilemanager}
			defaultInterface, defaultGateway, defaultSubnet, err := r.DefaultRoute()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.defaultInterface, defaultInterface)
			assert.Equal(t, tc.defaultGateway, defaultGateway)
			assert.Equal(t, tc.defaultSubnet, defaultSubnet)
		})
	}
}

func Test_routeExists(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		subnet  net.IPNet
		data    []byte
		readErr error
		exists  bool
		err     error
	}{
		"no data": {},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"parse error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x`),
			err: fmt.Errorf("line 1 in /proc/net/route: line \"eth0   x\": not enough fields"),
		},
		"not existing": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 128},
			},
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    0002A8C0        0100000A        0003    0       0       0       00FFFFFF        0       0       0`),
		},
		"existing": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    0002A8C0        0100000A        0003    0       0       0       00FFFFFF        0       0       0`),
			exists: true,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockFilemanager := NewMockFileManager(mockCtrl)

			mockFilemanager.EXPECT().ReadFile(string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Times(1)
			r := &routing{fileManager: mockFilemanager}
			exists, err := r.routeExists(tc.subnet)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.exists, exists)
		})
	}
}
