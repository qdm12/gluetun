package routing

import (
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/files/mock_files"
	"github.com/qdm12/golibs/logging/mock_logging"
)

const exampleRouteData = `Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT                                                       
tun0    00000000        050A030A        0003    0       0       0       00000080        0       0       0                                                                               
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0                                                                               
tun0    010A030A        050A030A        0007    0       0       0       FFFFFFFF        0       0       0                                                                               
tun0    050A030A        00000000        0005    0       0       0       FFFFFFFF        0       0       0                                                                               
eth0    42196956        010011AC        0007    0       0       0       FFFFFFFF        0       0       0                                                                               
tun0    00000080        050A030A        0003    0       0       0       00000080        0       0       0                                                                               
eth0    000011AC        00000000        0001    0       0       0       0000FFFF        0       0       0
`

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
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
`),
			entries: []routingEntry{},
		},
		"legend and single line": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   0002A8C0  0100000A  0003   0 0 0  00FFFFFF   0 0  0
`),
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
eth0   0002A8C0  0100000A  0002   0 0 0  00FFFFFF   0 0  0
`),
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
eth0   x  0100000A  0003   0 0 0  00FFFFFF   0 0  0
`),
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
		err              error
	}{
		"no data": {
			err: fmt.Errorf("not enough entries (0) found in %s", constants.NetRoute)},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error")},
		"parse error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x
`),
			err: fmt.Errorf("line 1 in /proc/net/route: line \"eth0   x\": not enough fields")},
		"single entry": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    00000000        050A090A        0003    0       0       0       00000080        0       0       0
`),
			err: fmt.Errorf("not enough entries (1) found in %s", constants.NetRoute)},
		"success": {
			data:             []byte(exampleRouteData),
			defaultInterface: "eth0",
			defaultGateway:   net.IP{172, 17, 0, 1},
		},
		"not found": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT        
eth0    00000000        010011AC        0003    0       0       0       10000000        0       0       0
eth0    000011AC        00000000        0001    0       0       0       0000FFFF        0       0       0
`),
			err: fmt.Errorf("cannot find default route"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			filemanager := mock_files.NewMockFileManager(mockCtrl)

			filemanager.EXPECT().ReadFile(string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Times(1)
			if tc.err == nil {
				logger.EXPECT().Info(
					"default route found: interface %s, gateway %s",
					tc.defaultInterface, tc.defaultGateway.String(),
				).Times(1)
			}
			r := &routing{logger: logger, fileManager: filemanager}
			defaultInterface, defaultGateway, err := r.DefaultRoute()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.defaultInterface, defaultInterface)
			assert.Equal(t, tc.defaultGateway, defaultGateway)
		})
	}
}

func Test_LocalSubnet(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data        []byte
		readErr     error
		localSubnet net.IPNet
		err         error
	}{
		"no data": {
			err: fmt.Errorf("not enough entries (0) found in %s", constants.NetRoute)},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error")},
		"parse error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x
`),
			err: fmt.Errorf("line 1 in /proc/net/route: line \"eth0   x\": not enough fields")},
		"single entry": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT        
eth0    00000000        050A090A        0003    0       0       0       00000080        0       0       0
`),
			err: fmt.Errorf("not enough entries (1) found in %s", constants.NetRoute)},
		"success": {
			data: []byte(exampleRouteData),
			localSubnet: net.IPNet{
				IP:   net.IP{172, 17, 0, 0},
				Mask: net.IPMask{255, 255, 0, 0},
			},
		},
		"not found": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT        
eth0    00000000        010011AC        0003    0       0       0       00000000        0       0       0
eth0    000011AC        10000000        0001    0       0       0       0000FFFF        0       0       0
`),
			err: fmt.Errorf("cannot find local subnet route"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			filemanager := mock_files.NewMockFileManager(mockCtrl)

			filemanager.EXPECT().ReadFile(string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Times(1)
			if tc.err == nil {
				logger.EXPECT().Info("local subnet found: %s", tc.localSubnet.String()).Times(1)
			}
			r := &routing{logger: logger, fileManager: filemanager}
			localSubnet, err := r.LocalSubnet()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.localSubnet, localSubnet)
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
			err:     fmt.Errorf("cannot check route existence: error"),
		},
		"parse error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x
`),
			err: fmt.Errorf("cannot check route existence: line 1 in /proc/net/route: line \"eth0   x\": not enough fields"),
		},
		"not existing": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 128},
			},
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    0002A8C0        0100000A        0003    0       0       0       00FFFFFF        0       0       0
`),
		},
		"existing": {
			subnet: net.IPNet{
				IP:   net.IP{192, 168, 2, 0},
				Mask: net.IPMask{255, 255, 255, 0},
			},
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0    0002A8C0        0100000A        0003    0       0       0       00FFFFFF        0       0       0
`),
			exists: true,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			filemanager := mock_files.NewMockFileManager(mockCtrl)
			filemanager.EXPECT().ReadFile(string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Times(1)
			r := &routing{fileManager: filemanager}
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

func Test_VPNGatewayIP(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		defaultInterface string
		data             []byte
		readErr          error
		ip               net.IP
		err              error
	}{
		"no data": {
			err: fmt.Errorf("cannot find VPN gateway IP address from ip routes"),
		},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("cannot find VPN gateway IP address: error"),
		},
		"parse error": {
			data: []byte(`Iface   Destination     Gateway         Flags   RefCnt  Use     Metric  Mask            MTU     Window  IRTT
eth0   x
`),
			err: fmt.Errorf("cannot find VPN gateway IP address: line 1 in /proc/net/route: line \"eth0   x\": not enough fields"),
		},
		"found eth0": {
			defaultInterface: "eth0",
			data:             []byte(exampleRouteData),
			ip:               net.IP{86, 105, 25, 66},
		},
		"not found tun0": {
			defaultInterface: "tun0",
			data:             []byte(exampleRouteData),
			err:              fmt.Errorf("cannot find VPN gateway IP address from ip routes"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			filemanager := mock_files.NewMockFileManager(mockCtrl)
			filemanager.EXPECT().ReadFile(string(constants.NetRoute)).
				Return(tc.data, tc.readErr).Times(1)
			r := &routing{fileManager: filemanager}
			ip, err := r.VPNGatewayIP(tc.defaultInterface)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.ip, ip)
		})
	}
}
