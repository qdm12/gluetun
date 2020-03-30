package shadowsocks

import (
	"fmt"
	"testing"

	filesMocks "github.com/qdm12/golibs/files/mocks"
	loggingMocks "github.com/qdm12/golibs/logging/mocks"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
			logger := &loggingMocks.Logger{}
			logger.On("Info", "%s: generating configuration file", logPrefix).Once()
			fileManager := &filesMocks.FileManager{}
			fileManager.On("WriteToFile",
				string(constants.ShadowsocksConf),
				[]byte(`{"server":"0.0.0.0","user":"nonrootuser","method":"chacha20-ietf-poly1305","timeout":30,"fast_open":false,"mode":"tcp_and_udp","port_password":{"2000":"abcde"},"workers":2,"interface":"tun","nameserver":"127.0.0.1"}`),
				mock.AnythingOfType("files.WriteOptionSetter"),
				mock.AnythingOfType("files.WriteOptionSetter"),
			).
				Return(tc.writeErr).Once()
			c := &configurator{logger: logger, fileManager: fileManager}
			err := c.MakeConf(2000, "abcde", "chacha20-ietf-poly1305", 1000, 1001)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			logger.AssertExpectations(t)
			fileManager.AssertExpectations(t)
		})
	}
}
