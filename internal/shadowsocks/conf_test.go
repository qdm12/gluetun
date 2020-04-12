package shadowsocks

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/files/mock_files"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_generateConf(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		port     uint16
		password string
		data     []byte
	}{
		"no data": {
			data: []byte(`{"server":"0.0.0.0","user":"nonrootuser","method":"chacha20-ietf-poly1305","timeout":30,"fast_open":false,"mode":"tcp_and_udp","port_password":{"0":""},"workers":2,"interface":"tun","nameserver":"127.0.0.1"}`),
		},
		"data": {
			port:     2000,
			password: "abcde",
			data:     []byte(`{"server":"0.0.0.0","user":"nonrootuser","method":"chacha20-ietf-poly1305","timeout":30,"fast_open":false,"mode":"tcp_and_udp","port_password":{"2000":"abcde"},"workers":2,"interface":"tun","nameserver":"127.0.0.1"}`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			data := generateConf(tc.port, tc.password, "chacha20-ietf-poly1305")
			assert.Equal(t, tc.data, data)
		})
	}
}

func Test_MakeConf(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		writeErr error
		err      error
	}{
		"no write error": {},
		"write error": {
			writeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("generating configuration file").Times(1)
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			fileManager.EXPECT().WriteToFile(
				string(constants.ShadowsocksConf),
				[]byte(`{"server":"0.0.0.0","user":"nonrootuser","method":"chacha20-ietf-poly1305","timeout":30,"fast_open":false,"mode":"tcp_and_udp","port_password":{"2000":"abcde"},"workers":2,"interface":"tun","nameserver":"127.0.0.1"}`),
				gomock.AssignableToTypeOf(files.Ownership(0, 0)),
				gomock.AssignableToTypeOf(files.Ownership(0, 0)),
			).Return(tc.writeErr).Times(1)
			c := &configurator{logger: logger, fileManager: fileManager}
			err := c.MakeConf(2000, "abcde", "chacha20-ietf-poly1305", 1000, 1001)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
