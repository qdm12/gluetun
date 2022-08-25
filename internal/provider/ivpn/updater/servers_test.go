package updater

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/provider/common"
	"github.com/qdm12/gluetun/internal/updater/resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Updater_GetServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		// Inputs
		minServers int

		// Mocks
		warnerBuilder func(ctrl *gomock.Controller) common.Warner

		// From API
		responseBody   string
		responseStatus int

		// Resolution
		expectResolve   bool
		resolveSettings resolver.ParallelSettings
		hostToIPs       map[string][]net.IP
		resolveWarnings []string
		resolveErr      error

		// Output
		servers []models.Server
		err     error
	}{
		"http response error": {
			warnerBuilder:  func(ctrl *gomock.Controller) common.Warner { return nil },
			responseStatus: http.StatusNoContent,
			err:            errors.New("failed fetching API: HTTP status code not OK: 204 No Content"),
		},
		"resolve error": {
			warnerBuilder: func(ctrl *gomock.Controller) common.Warner {
				warner := common.NewMockWarner(ctrl)
				warner.EXPECT().Warn("resolve warning")
				return warner
			},
			responseBody: `{"servers":[
				{"hostnames":{"openvpn":"hosta"}}
			]}`,
			responseStatus: http.StatusOK,
			expectResolve:  true,
			resolveSettings: resolver.ParallelSettings{
				Hosts:        []string{"hosta"},
				MaxFailRatio: 0.1,
				Repeat: resolver.RepeatSettings{
					MaxDuration:     20 * time.Second,
					BetweenDuration: time.Second,
					MaxNoNew:        2,
					MaxFails:        2,
					SortIPs:         true,
				},
			},
			resolveWarnings: []string{"resolve warning"},
			resolveErr:      errors.New("dummy"),
			err:             errors.New("dummy"),
		},
		"not enough servers": {
			minServers:    2,
			warnerBuilder: func(ctrl *gomock.Controller) common.Warner { return nil },
			responseBody: `{"servers":[
				{"hostnames":{"openvpn":"hosta"}}
			]}`,
			responseStatus: http.StatusOK,
			err:            errors.New("not enough servers found: 1 and expected at least 2"),
		},
		"success": {
			minServers: 1,
			warnerBuilder: func(ctrl *gomock.Controller) common.Warner {
				warner := common.NewMockWarner(ctrl)
				warner.EXPECT().Warn("resolve warning")
				return warner
			},
			responseBody: `{"servers":[
				{"country":"Country1","city":"City A","hostnames":{"openvpn":"hosta"}},
				{"country":"Country2","city":"City B","hostnames":{"openvpn":"hostb"},"wg_public_key":"xyz"},
				{"country":"Country3","city":"City C","hostnames":{"wireguard":"hostc"},"wg_public_key":"xyz"}
			]}`,
			responseStatus: http.StatusOK,
			expectResolve:  true,
			resolveSettings: resolver.ParallelSettings{
				Hosts:        []string{"hosta", "hostb", "hostc"},
				MaxFailRatio: 0.1,
				Repeat: resolver.RepeatSettings{
					MaxDuration:     20 * time.Second,
					BetweenDuration: time.Second,
					MaxNoNew:        2,
					MaxFails:        2,
					SortIPs:         true,
				},
			},
			hostToIPs: map[string][]net.IP{
				"hosta": {{1, 1, 1, 1}, {2, 2, 2, 2}},
				"hostb": {{3, 3, 3, 3}, {4, 4, 4, 4}},
				"hostc": {{5, 5, 5, 5}, {6, 6, 6, 6}},
			},
			resolveWarnings: []string{"resolve warning"},
			servers: []models.Server{
				{VPN: vpn.OpenVPN, Country: "Country1",
					City: "City A", Hostname: "hosta", TCP: true, UDP: true,
					IPs: []net.IP{{1, 1, 1, 1}, {2, 2, 2, 2}}},
				{VPN: vpn.OpenVPN, Country: "Country2",
					City: "City B", Hostname: "hostb", TCP: true, UDP: true,
					IPs: []net.IP{{3, 3, 3, 3}, {4, 4, 4, 4}}},
				{VPN: vpn.Wireguard,
					Country: "Country3", City: "City C",
					Hostname: "hostc",
					WgPubKey: "xyz",
					IPs:      []net.IP{{5, 5, 5, 5}, {6, 6, 6, 6}}},
			},
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			ctx := context.Background()

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, r.URL.String(), "https://api.ivpn.net/v4/servers/stats")
					return &http.Response{
						StatusCode: testCase.responseStatus,
						Status:     http.StatusText(testCase.responseStatus),
						Body:       io.NopCloser(strings.NewReader(testCase.responseBody)),
					}, nil
				}),
			}

			parallelResolver := common.NewMockParallelResolver(ctrl)
			if testCase.expectResolve {
				parallelResolver.EXPECT().Resolve(ctx, testCase.resolveSettings).
					Return(testCase.hostToIPs, testCase.resolveWarnings, testCase.resolveErr)
			}

			warner := testCase.warnerBuilder(ctrl)

			updater := &Updater{
				client:           client,
				parallelResolver: parallelResolver,
				warner:           warner,
			}

			servers, err := updater.FetchServers(ctx, testCase.minServers)

			assert.Equal(t, testCase.servers, servers)
			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
