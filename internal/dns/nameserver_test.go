package dns

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/os"
	"github.com/qdm12/gluetun/internal/os/mock_os"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UseDNSSystemWide(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		data        []byte
		writtenData string
		openErr     error
		readErr     error
		writeErr    error
		closeErr    error
		err         error
	}{
		"no data": {
			writtenData: "nameserver 127.0.0.1\n",
		},
		"open error": {
			openErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"write error": {
			writtenData: "nameserver 127.0.0.1\n",
			writeErr:    fmt.Errorf("error"),
			err:         fmt.Errorf("error"),
		},
		"lines without nameserver": {
			data:        []byte("abc\ndef\n"),
			writtenData: "abc\ndef\nnameserver 127.0.0.1\n",
		},
		"lines with nameserver": {
			data:        []byte("abc\nnameserver abc def\ndef\n"),
			writtenData: "abc\nnameserver 127.0.0.1\ndef\n",
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)

			file := mock_os.NewMockFile(mockCtrl)
			if tc.openErr == nil {
				firstReadCall := file.EXPECT().
					Read(gomock.AssignableToTypeOf([]byte{})).
					DoAndReturn(func(b []byte) (int, error) {
						copy(b, tc.data)
						return len(tc.data), nil
					})
				readErr := tc.readErr
				if readErr == nil {
					readErr = io.EOF
				}
				finalReadCall := file.EXPECT().
					Read(gomock.AssignableToTypeOf([]byte{})).
					Return(0, readErr).After(firstReadCall)
				if tc.readErr == nil {
					writeCall := file.EXPECT().WriteString(tc.writtenData).
						Return(0, tc.writeErr).After(finalReadCall)
					file.EXPECT().Close().Return(tc.closeErr).After(writeCall)
				} else {
					file.EXPECT().Close().Return(tc.closeErr).After(finalReadCall)
				}
			}

			openFile := func(name string, flag int, perm os.FileMode) (os.File, error) {
				assert.Equal(t, string(constants.ResolvConf), name)
				assert.Equal(t, os.O_RDWR|os.O_TRUNC, flag)
				assert.Equal(t, os.FileMode(0644), perm)
				return file, tc.openErr
			}

			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("using DNS address %s system wide", "127.0.0.1")
			c := &configurator{
				openFile: openFile,
				logger:   logger,
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
