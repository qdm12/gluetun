package dns

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/os"
	"github.com/qdm12/gluetun/internal/os/mock_os"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/golibs/network/mock_network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_downloadAndSave(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		content   []byte
		status    int
		clientErr error
		openErr   error
		writeErr  error
		chownErr  error
		closeErr  error
		err       error
	}{
		"no data": {
			status: http.StatusOK,
		},
		"bad status": {
			status: http.StatusBadRequest,
			err:    fmt.Errorf("HTTP status code is 400 for https://raw.githubusercontent.com/qdm12/files/master/named.root.updated"), //nolint:lll
		},
		"client error": {
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf("error"),
		},
		"open error": {
			status:  http.StatusOK,
			openErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"write error": {
			status:   http.StatusOK,
			writeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"chown error": {
			status:   http.StatusOK,
			chownErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"close error": {
			status:   http.StatusOK,
			closeErr: fmt.Errorf("error"),
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

			ctx := context.Background()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("downloading %s from %s", "root hints", string(constants.NamedRootURL))
			client := mock_network.NewMockClient(mockCtrl)
			client.EXPECT().Get(ctx, string(constants.NamedRootURL)).
				Return(tc.content, tc.status, tc.clientErr)

			openFile := func(name string, flag int, perm os.FileMode) (os.File, error) {
				return nil, nil
			}

			if tc.clientErr == nil && tc.status == http.StatusOK {
				file := mock_os.NewMockFile(mockCtrl)
				if tc.openErr == nil {
					writeCall := file.EXPECT().Write(tc.content).
						Return(0, tc.writeErr)
					if tc.writeErr != nil {
						file.EXPECT().Close().Return(tc.closeErr).After(writeCall)
					} else {
						chownCall := file.EXPECT().Chown(1000, 1000).Return(tc.chownErr).After(writeCall)
						file.EXPECT().Close().Return(tc.closeErr).After(chownCall)
					}
				}

				openFile = func(name string, flag int, perm os.FileMode) (os.File, error) {
					assert.Equal(t, string(constants.RootHints), name)
					assert.Equal(t, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, flag)
					assert.Equal(t, os.FileMode(0400), perm)
					return file, tc.openErr
				}
			}

			c := &configurator{
				logger:   logger,
				client:   client,
				openFile: openFile,
			}

			err := c.downloadAndSave(ctx, "root hints",
				string(constants.NamedRootURL), string(constants.RootHints),
				1000, 1000)

			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func Test_DownloadRootHints(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)

	ctx := context.Background()
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Info("downloading %s from %s", "root hints", string(constants.NamedRootURL))
	client := mock_network.NewMockClient(mockCtrl)
	client.EXPECT().Get(ctx, string(constants.NamedRootURL)).
		Return(nil, http.StatusOK, errors.New("test"))

	c := &configurator{
		logger: logger,
		client: client,
	}

	err := c.DownloadRootHints(ctx, 1000, 1000)
	require.Error(t, err)
	assert.Equal(t, "test", err.Error())
}

func Test_DownloadRootKey(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)

	ctx := context.Background()
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Info("downloading %s from %s", "root key", string(constants.RootKeyURL))
	client := mock_network.NewMockClient(mockCtrl)
	client.EXPECT().Get(ctx, string(constants.RootKeyURL)).
		Return(nil, http.StatusOK, errors.New("test"))

	c := &configurator{
		logger: logger,
		client: client,
	}

	err := c.DownloadRootKey(ctx, 1000, 1000)
	require.Error(t, err)
	assert.Equal(t, "test", err.Error())
}
