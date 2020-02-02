package dns

import (
	"fmt"
	"testing"

	"github.com/qdm12/golibs/files/mocks"
	"github.com/qdm12/private-internet-access-docker/internal/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO check empty lines write 0 byte in golibs/files

func Test_SetLocalNameserver(t *testing.T) {
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
			fileManager := &mocks.FileManager{}
			fileManager.On("ReadFile", string(constants.ResolvConf)).
				Return(tc.data, tc.readErr).Once()
			fileManager.On("WriteToFile", string(constants.ResolvConf), tc.writtenData).
				Return(tc.writeErr).Once()
			c := &configurator{
				fileManager: fileManager,
			}
			err := c.SetLocalNameserver()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
