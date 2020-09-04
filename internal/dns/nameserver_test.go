package dns

import (
	"fmt"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/files/mock_files"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UseDNSSystemWide(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data        []byte
		writtenData []byte
		readErr     error
		writeErr    error
		err         error
	}{
		"no data": {
			writtenData: []byte("nameserver 127.0.0.1"),
		},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"write error": {
			writtenData: []byte("nameserver 127.0.0.1"),
			writeErr:    fmt.Errorf("error"),
			err:         fmt.Errorf("error"),
		},
		"lines without nameserver": {
			data:        []byte("abc\ndef\n"),
			writtenData: []byte("abc\ndef\nnameserver 127.0.0.1"),
		},
		"lines with nameserver": {
			data:        []byte("abc\nnameserver abc def\ndef\n"),
			writtenData: []byte("abc\nnameserver 127.0.0.1\ndef"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			fileManager.EXPECT().ReadFile(string(constants.ResolvConf)).
				Return(tc.data, tc.readErr).Times(1)
			if tc.readErr == nil {
				fileManager.EXPECT().WriteToFile(string(constants.ResolvConf), tc.writtenData).
					Return(tc.writeErr).Times(1)
			}
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("using DNS address %s system wide", "127.0.0.1").Times(1)
			c := &configurator{
				fileManager: fileManager,
				logger:      logger,
			}
			err := c.UseDNSSystemWide(net.IP{127, 0, 0, 1}, false)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
