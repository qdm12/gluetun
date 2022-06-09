package ivpn

import (
	"context"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/gluetun/internal/updater/resolver/mock_resolver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Updater_GetServers(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		// Inputs
		minServers int

		// Mocks
		warnerBuilder func(ctrl *gomock.Controller) Warner

		// From API
		responseBody   string
		responseStatus int

		// Resolution
		expectResolve   bool
		hostsToResolve  []string
		hostToIPs       map[string][]net.IP
		resolveWarnings []string
		resolveErr      error

		// Output
		servers []models.Server
		err     error
	}{
		"http response error": {
			warnerBuilder:  func(ctrl *gomock.Controller) Warner { return nil },
			responseStatus: http.StatusNoContent,
			err:            errors.New("failed fetching API: HTTP status code not OK: 204 No Content"),
		},
		"resolve error": {
			warnerBuilder: func(ctrl *gomock.Controller) Warner {
				warner := NewMockWarner(ctrl)
				warner.EXPECT().Warn("resolve warning")
				return warner
			},
			responseBody: `{"servers":[
				{"hostnames":{"openvpn":"hosta"}}
			]}`,
			responseStatus:  http.StatusOK,
			expectResolve:   true,
			hostsToResolve:  []string{"hosta"},
			resolveWarnings: []string{"resolve warning"},
			resolveErr:      errors.New("dummy"),
			err:             errors.New("dummy"),
		},
		"not enough servers": {
			minServers:    2,
			warnerBuilder: func(ctrl *gomock.Controller) Warner { return nil },
			responseBody: `{"servers":[
				{"hostnames":{"openvpn":"hosta"}}
			]}`,
			responseStatus: http.StatusOK,
			err:            errors.New("not enough servers found: 1 and expected at least 2"),
		},
		"success": {
			minServers: 1,
			warnerBuilder: func(ctrl *gomock.Controller) Warner {
				warner := NewMockWarner(ctrl)
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
			hostsToResolve: []string{"hosta", "hostb", "hostc"},
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
					Hostname: "hostc", UDP: true,
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
						Body:       ioutil.NopCloser(strings.NewReader(testCase.responseBody)),
					}, nil
				}),
			}

			presolver := mock_resolver.NewMockParallel(ctrl)
			if testCase.expectResolve {
				presolver.EXPECT().Resolve(ctx, testCase.hostsToResolve, testCase.minServers).
					Return(testCase.hostToIPs, testCase.resolveWarnings, testCase.resolveErr)
			}

			warner := testCase.warnerBuilder(ctrl)

			updater := &Updater{
				client:    client,
				presolver: presolver,
				warner:    warner,
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
