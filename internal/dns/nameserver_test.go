package dns

import (
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/os/mock_os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_UseDNSSystemWide(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		ip             net.IP
		keepNameserver bool
		data           []byte
		writtenData    string
		openErr        error
		readErr        error
		writeErr       error
		closeErr       error
		err            error
	}{
		"no data": {
			ip:          net.IP{127, 0, 0, 1},
			writtenData: "nameserver 127.0.0.1\n",
		},
		"open error": {
			ip:      net.IP{127, 0, 0, 1},
			openErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"write error": {
			ip:          net.IP{127, 0, 0, 1},
			writtenData: "nameserver 127.0.0.1\n",
			writeErr:    fmt.Errorf("error"),
			err:         fmt.Errorf("error"),
		},
		"lines without nameserver": {
			ip:          net.IP{127, 0, 0, 1},
			data:        []byte("abc\ndef\n"),
			writtenData: "nameserver 127.0.0.1\nabc\ndef\n",
		},
		"lines with nameserver": {
			ip:          net.IP{127, 0, 0, 1},
			data:        []byte("abc\nnameserver abc def\ndef\n"),
			writtenData: "nameserver 127.0.0.1\nabc\ndef\n",
		},
		"keep nameserver": {
			ip:             net.IP{127, 0, 0, 1},
			keepNameserver: true,
			data:           []byte("abc\nnameserver abc def\ndef\n"),
			writtenData:    "nameserver 127.0.0.1\nabc\nnameserver abc def\ndef\n",
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
			logger.EXPECT().Info("using DNS address %s system wide", tc.ip.String())
			c := &configurator{
				openFile: openFile,
				logger:   logger,
			}
			err := c.UseDNSSystemWide(tc.ip, tc.keepNameserver)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
