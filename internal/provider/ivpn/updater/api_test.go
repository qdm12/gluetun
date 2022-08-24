package updater

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_fetchAPI(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		responseStatus int
		responseBody   io.ReadCloser
		data           apiData
		err            error
	}{
		"http response status not ok": {
			responseStatus: http.StatusNoContent,
			err:            errors.New("HTTP status code not OK: 204 No Content"),
		},
		"nil body": {
			responseStatus: http.StatusOK,
			err:            errors.New("failed unmarshaling response body: EOF"),
		},
		"no server": {
			responseStatus: http.StatusOK,
			responseBody:   io.NopCloser(strings.NewReader(`{}`)),
		},
		"success": {
			responseStatus: http.StatusOK,
			responseBody: io.NopCloser(strings.NewReader(`{"servers":[
				{"country":"Country1","city":"City A","isp":"xyz","is_active":true,"hostnames":{"openvpn":"hosta"}},
				{"country":"Country2","city":"City B","isp":"abc","is_active":false,"hostnames":{"openvpn":"hostb"}}
			]}`)),
			data: apiData{
				Servers: []apiServer{
					{
						Country:  "Country1",
						City:     "City A",
						IsActive: true,
						ISP:      "xyz",
						Hostnames: apiHostnames{
							OpenVPN: "hosta",
						},
					},
					{
						Country: "Country2",
						City:    "City B",
						ISP:     "abc",
						Hostnames: apiHostnames{
							OpenVPN: "hostb",
						},
					},
				},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, r.URL.String(), "https://api.ivpn.net/v4/servers/stats")
					return &http.Response{
						StatusCode: testCase.responseStatus,
						Status:     http.StatusText(testCase.responseStatus),
						Body:       testCase.responseBody,
					}, nil
				}),
			}

			data, err := fetchAPI(ctx, client)

			assert.Equal(t, testCase.data, data)
			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
