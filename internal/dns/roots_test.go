package dns

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/files/mock_files"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/golibs/network/mock_network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/gluetun/internal/constants"
)

func Test_DownloadRootHints(t *testing.T) { //nolint:dupl
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("downloading root hints from %s", constants.NamedRootURL).Times(1)
			client := mock_network.NewMockClient(mockCtrl)
			client.EXPECT().GetContent(string(constants.NamedRootURL)).
				Return(tc.content, tc.status, tc.clientErr).Times(1)
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			if tc.clientErr == nil && tc.status == http.StatusOK {
				fileManager.EXPECT().WriteToFile(
					string(constants.RootHints),
					tc.content,
					gomock.AssignableToTypeOf(files.Ownership(0, 0)),
					gomock.AssignableToTypeOf(files.Ownership(0, 0))).
					Return(tc.writeErr).Times(1)
			}
			c := &configurator{logger: logger, client: client, fileManager: fileManager}
			err := c.DownloadRootHints(1000, 1000)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_DownloadRootKey(t *testing.T) { //nolint:dupl
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("downloading root key from %s", constants.RootKeyURL).Times(1)
			client := mock_network.NewMockClient(mockCtrl)
			client.EXPECT().GetContent(string(constants.RootKeyURL)).
				Return(tc.content, tc.status, tc.clientErr).Times(1)
			fileManager := mock_files.NewMockFileManager(mockCtrl)
			if tc.clientErr == nil && tc.status == http.StatusOK {
				fileManager.EXPECT().WriteToFile(
					string(constants.RootKey),
					tc.content,
					gomock.AssignableToTypeOf(files.Ownership(0, 0)),
					gomock.AssignableToTypeOf(files.Ownership(0, 0)),
				).Return(tc.writeErr).Times(1)
			}
			c := &configurator{logger: logger, client: client, fileManager: fileManager}
			err := c.DownloadRootKey(1000, 1001)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
