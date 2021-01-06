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
		firstOpenErr   error
		readErr        error
		firstCloseErr  error
		secondOpenErr  error
		writtenData    string
		writeErr       error
		secondCloseErr error
		err            error
	}{
		"no data": {
			ip:          net.IP{127, 0, 0, 1},
			writtenData: "nameserver 127.0.0.1\n",
		},
		"first open error": {
			ip:           net.IP{127, 0, 0, 1},
			firstOpenErr: fmt.Errorf("error"),
			err:          fmt.Errorf("error"),
		},
		"read error": {
			readErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"first close error": {
			firstCloseErr: fmt.Errorf("error"),
			err:           fmt.Errorf("error"),
		},
		"second open error": {
			ip:            net.IP{127, 0, 0, 1},
			secondOpenErr: fmt.Errorf("error"),
			err:           fmt.Errorf("error"),
		},
		"write error": {
			ip:          net.IP{127, 0, 0, 1},
			writtenData: "nameserver 127.0.0.1\n",
			writeErr:    fmt.Errorf("error"),
			err:         fmt.Errorf("error"),
		},
		"second close error": {
			ip:             net.IP{127, 0, 0, 1},
			writtenData:    "nameserver 127.0.0.1\n",
			secondCloseErr: fmt.Errorf("error"),
			err:            fmt.Errorf("error"),
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

			type fileCall struct {
				path string
				flag int
				perm os.FileMode
				file os.File
				err  error
			}

			var fileCalls []fileCall

			readOnlyFile := mock_os.NewMockFile(mockCtrl)

			if tc.firstOpenErr == nil {
				firstReadCall := readOnlyFile.EXPECT().
					Read(gomock.AssignableToTypeOf([]byte{})).
					DoAndReturn(func(b []byte) (int, error) {
						copy(b, tc.data)
						return len(tc.data), nil
					})
				readErr := tc.readErr
				if readErr == nil {
					readErr = io.EOF
				}
				finalReadCall := readOnlyFile.EXPECT().
					Read(gomock.AssignableToTypeOf([]byte{})).
					Return(0, readErr).After(firstReadCall)
				readOnlyFile.EXPECT().Close().
					Return(tc.firstCloseErr).
					After(finalReadCall)
			}

			fileCalls = append(fileCalls, fileCall{
				path: string(constants.ResolvConf),
				flag: os.O_RDONLY,
				perm: 0,
				file: readOnlyFile,
				err:  tc.firstOpenErr,
			}) // always return readOnlyFile

			if tc.firstOpenErr == nil && tc.readErr == nil && tc.firstCloseErr == nil {
				writeOnlyFile := mock_os.NewMockFile(mockCtrl)
				if tc.secondOpenErr == nil {
					writeCall := writeOnlyFile.EXPECT().
						WriteString(tc.writtenData).
						Return(0, tc.writeErr)
					writeOnlyFile.EXPECT().
						Close().
						Return(tc.secondCloseErr).
						After(writeCall)
				}
				fileCalls = append(fileCalls, fileCall{
					path: string(constants.ResolvConf),
					flag: os.O_WRONLY | os.O_TRUNC,
					perm: os.FileMode(0644),
					file: writeOnlyFile,
					err:  tc.secondOpenErr,
				})
			}

			fileCallIndex := 0
			openFile := func(name string, flag int, perm os.FileMode) (os.File, error) {
				fileCall := fileCalls[fileCallIndex]
				fileCallIndex++
				assert.Equal(t, fileCall.path, name)
				assert.Equal(t, fileCall.flag, flag)
				assert.Equal(t, fileCall.perm, perm)
				return fileCall.file, fileCall.err
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
