package dns

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/qdm12/golibs/os"
	"github.com/qdm12/golibs/os/mock_os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_downloadAndSave(t *testing.T) {
	t.Parallel()
	const defaultURL = "https://test.com"
	tests := map[string]struct {
		url       string // to trigger a new request error
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
			url:    defaultURL,
			status: http.StatusOK,
		},
		"bad status": {
			url:    defaultURL,
			status: http.StatusBadRequest,
			err:    fmt.Errorf("bad HTTP status from %s: Bad Request", defaultURL),
		},
		"client error": {
			url:       defaultURL,
			clientErr: fmt.Errorf("error"),
			err:       fmt.Errorf("Get %q: error", defaultURL),
		},
		"open error": {
			url:     defaultURL,
			status:  http.StatusOK,
			openErr: fmt.Errorf("error"),
			err:     fmt.Errorf("error"),
		},
		"chown error": {
			url:      defaultURL,
			status:   http.StatusOK,
			chownErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"close error": {
			url:      defaultURL,
			status:   http.StatusOK,
			closeErr: fmt.Errorf("error"),
			err:      fmt.Errorf("error"),
		},
		"data": {
			url:     defaultURL,
			content: []byte("content"),
			status:  http.StatusOK,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockCtrl := gomock.NewController(t)

			ctx := context.Background()
			logger := mock_logging.NewMockLogger(mockCtrl)
			logger.EXPECT().Info("downloading %s from %s", "root hints", tc.url)

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, tc.url, r.URL.String())
					if tc.clientErr != nil {
						return nil, tc.clientErr
					}
					return &http.Response{
						StatusCode: tc.status,
						Status:     http.StatusText(tc.status),
						Body:       ioutil.NopCloser(bytes.NewReader(tc.content)),
					}, nil
				}),
			}

			openFile := func(name string, flag int, perm os.FileMode) (os.File, error) {
				return nil, nil
			}

			const filepath = "/test"

			if tc.clientErr == nil && tc.status == http.StatusOK {
				file := mock_os.NewMockFile(mockCtrl)
				if tc.openErr == nil {
					if len(tc.content) > 0 {
						file.EXPECT().
							Write(tc.content).
							Return(len(tc.content), tc.writeErr)
					}
					file.EXPECT().
						Close().
						Return(tc.closeErr)
					file.EXPECT().
						Chown(1000, 1000).
						Return(tc.chownErr)
				}

				openFile = func(name string, flag int, perm os.FileMode) (os.File, error) {
					assert.Equal(t, filepath, name)
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
				tc.url, filepath,
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

	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, string(constants.NamedRootURL), r.URL.String())
			return nil, errors.New("test")
		}),
	}

	c := &configurator{
		logger: logger,
		client: client,
	}

	err := c.DownloadRootHints(ctx, 1000, 1000)
	require.Error(t, err)
	assert.Equal(t, `Get "https://raw.githubusercontent.com/qdm12/files/master/named.root.updated": test`, err.Error())
}

func Test_DownloadRootKey(t *testing.T) {
	t.Parallel()
	mockCtrl := gomock.NewController(t)

	ctx := context.Background()
	logger := mock_logging.NewMockLogger(mockCtrl)
	logger.EXPECT().Info("downloading %s from %s", "root key", string(constants.RootKeyURL))

	client := &http.Client{
		Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, string(constants.RootKeyURL), r.URL.String())
			return nil, errors.New("test")
		}),
	}

	c := &configurator{
		logger: logger,
		client: client,
	}

	err := c.DownloadRootKey(ctx, 1000, 1000)
	require.Error(t, err)
	assert.Equal(t, `Get "https://raw.githubusercontent.com/qdm12/files/master/root.key.updated": test`, err.Error())
}
