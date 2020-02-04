package dns

import (
	"fmt"
	"net/http"
	"testing"

	filesMocks "github.com/qdm12/golibs/files/mocks"
	loggingMocks "github.com/qdm12/golibs/logging/mocks"
	networkMocks "github.com/qdm12/golibs/network/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/private-internet-access-docker/internal/constants"
)

func Test_DownloadRootHints(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		content   []byte
		status    int
		clientErr error
		writeErr  error
		err       error
	}{
		"no data": {
			status: http.StatusOK,
		},
		"bad status": {
			status: http.StatusBadRequest,
			err:    fmt.Errorf("HTTP status code is 400 for https://raw.githubusercontent.com/qdm12/files/master/named.root.updated"),
		},
		"client error": {
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf("error"),
		},
		"write error": {
			status:   http.StatusOK,
			writeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"data": {
			content: []byte("content"),
			status:  http.StatusOK,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := &loggingMocks.Logger{}
			logger.On("Info", "%s: downloading root hints from %s", logPrefix, constants.NamedRootURL).Once()
			client := &networkMocks.Client{}
			client.On("GetContent", string(constants.NamedRootURL)).
				Return(tc.content, tc.status, tc.clientErr).Once()
			fileManager := &filesMocks.FileManager{}
			if tc.clientErr == nil && tc.status == http.StatusOK {
				fileManager.On("WriteToFile", string(constants.RootHints), tc.content).Return(tc.writeErr).Once()
			}
			c := &configurator{logger: logger, client: client, fileManager: fileManager}
			err := c.DownloadRootHints()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			logger.AssertExpectations(t)
			client.AssertExpectations(t)
			fileManager.AssertExpectations(t)
		})
	}
}

func Test_DownloadRootKey(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		content   []byte
		status    int
		clientErr error
		writeErr  error
		err       error
	}{
		"no data": {
			status: http.StatusOK,
		},
		"bad status": {
			status: http.StatusBadRequest,
			err:    fmt.Errorf("HTTP status code is 400 for https://raw.githubusercontent.com/qdm12/files/master/root.key.updated"),
		},
		"client error": {
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf("error"),
		},
		"write error": {
			status:   http.StatusOK,
			writeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"data": {
			content: []byte("content"),
			status:  http.StatusOK,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			logger := &loggingMocks.Logger{}
			logger.On("Info", "%s: downloading root key from %s", logPrefix, constants.RootKeyURL).Once()
			client := &networkMocks.Client{}
			client.On("GetContent", string(constants.RootKeyURL)).
				Return(tc.content, tc.status, tc.clientErr).Once()
			fileManager := &filesMocks.FileManager{}
			if tc.clientErr == nil && tc.status == http.StatusOK {
				fileManager.On("WriteToFile", string(constants.RootKey), tc.content).Return(tc.writeErr).Once()
			}
			c := &configurator{logger: logger, client: client, fileManager: fileManager}
			err := c.DownloadRootKey()
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			logger.AssertExpectations(t)
			client.AssertExpectations(t)
			fileManager.AssertExpectations(t)
		})
	}
}
