package updater

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

type roundTripFunc func(r *http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func Test_fetchServers(t *testing.T) {
	t.Parallel()

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	testCases := map[string]struct {
		ctx            context.Context
		responseStatus int
		responseBody   io.ReadCloser
		hostToData     map[string]serverData
		errWrapped     error
		errMessage     string
	}{
		"context canceled": {
			ctx:        canceledCtx,
			errWrapped: context.Canceled,
			errMessage: `fetching HTML code: Get "https://www.slickvpn.com/locations/": context canceled`,
		},
		"success": {
			ctx:            context.Background(),
			responseStatus: http.StatusOK,
			//nolint:lll
			responseBody: io.NopCloser(strings.NewReader(`
			<div>
				<table id="location-table">
					<tbody>
						<tr>
							<td>South America</td>
							<td>Chile</td>
							<td>Vina del Mar</td>
							<td> <a
									href="https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.kna1.slickvpn.com.ovpn">OpenVPN</a>
								| <a
									href="https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/viscosity/gw1.kna1.slickvpn.com_viscosity.ovpn">Viscosity</a>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
			`)),
			hostToData: map[string]serverData{
				"gw1.kna1.slickvpn.com": {
					ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.kna1.slickvpn.com.ovpn", //nolint:lll
					country: "Chile",
					region:  "South America",
					city:    "Vina del Mar",
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			client := &http.Client{
				Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodGet, r.Method)
					assert.Equal(t, r.URL.String(), "https://www.slickvpn.com/locations/")

					ctxErr := r.Context().Err()
					if ctxErr != nil {
						return nil, ctxErr
					}

					return &http.Response{
						StatusCode: http.StatusOK,
						Status:     http.StatusText(testCase.responseStatus),
						Body:       testCase.responseBody,
					}, nil
				}),
			}

			hostToData, err := fetchServers(testCase.ctx, client)

			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
			assert.Equal(t, testCase.hostToData, hostToData)
		})
	}
}

func Test_parseHTML(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		rootNode   *html.Node
		hostToData map[string]serverData
		errWrapped error
		errMessage string
	}{
		"empty html": {
			rootNode:   parseTestHTML(t, ""),
			errWrapped: ErrLocationTableNotFound,
			errMessage: `HTML location table node not found: in HTML code: <html><head></head><body></body></html>`,
		},
		"test data": {
			rootNode: parseTestDataIndexHTML(t),
			//nolint:lll
			hostToData: map[string]serverData{
				"gw1.ams1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ams1.slickvpn.com.ovpn", country: "Netherlands", region: "Europe", city: "Amsterdam"},
				"gw1.ams2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ams2.slickvpn.com.ovpn", country: "Netherlands", region: "Europe", city: "Amsterdam"},
				"gw1.ams3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ams3.slickvpn.com.ovpn", country: "Netherlands", region: "Europe", city: "Amsterdam"},
				"gw1.ams4.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ams4.slickvpn.com.ovpn", country: "Netherlands", region: "Europe", city: "Amsterdam"},
				"gw1.arn1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.arn1.slickvpn.com.ovpn", country: "Sweden", region: "Europe", city: "Stockholm"},
				"gw1.arn3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.arn3.slickvpn.com.ovpn", country: "Sweden", region: "Europe", city: "Stockholm"},
				"gw1.ath1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ath1.slickvpn.com.ovpn", country: "Greece", region: "Europe", city: "Athens"},
				"gw1.atl1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.atl1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Atlanta"},
				"gw1.atl3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.atl3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Atlanta"},
				"gw1.beg1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.beg1.slickvpn.com.ovpn", country: "Serbia", region: "Europe", city: "Belgrade"},
				"gw1.bkk1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.bkk1.slickvpn.com.ovpn", country: "Thailand", region: "Asia", city: "Bangkok"},
				"gw1.blr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.blr1.slickvpn.com.ovpn", country: "India", region: "Asia", city: "Bangalore"},
				"gw1.bne1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.bne1.slickvpn.com.ovpn", country: "Australia", region: "Oceania", city: "Brisbane"},
				"gw1.bom1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.bom1.slickvpn.com.ovpn", country: "India", region: "Asia", city: "Mumbai"},
				"gw1.bos1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.bos1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Boston"},
				"gw1.bud1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.bud1.slickvpn.com.ovpn", country: "Hungary", region: "Europe", city: "Budapest"},
				"gw1.buf1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.buf1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Buffalo"},
				"gw1.buh2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.buh2.slickvpn.com.ovpn", country: "Romania", region: "Europe", city: "Bucharest"},
				"gw1.cdg1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.cdg1.slickvpn.com.ovpn", country: "France", region: "Europe", city: "Paris"},
				"gw1.cgk1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.cgk1.slickvpn.com.ovpn", country: "Indonesia", region: "Asia", city: "Jakarta"},
				"gw1.cmh1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.cmh1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Columbus"},
				"gw1.cph1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.cph1.slickvpn.com.ovpn", country: "Denmark", region: "Europe", city: "Copenhagen"},
				"gw1.cvt1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.cvt1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "Coventry"},
				"gw1.dbq1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.dbq1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Dubuque"},
				"gw1.den1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.den1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Denver"},
				"gw1.dfw2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.dfw2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Dallas"},
				"gw1.dfw3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.dfw3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Dallas"},
				"gw1.dub1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.dub1.slickvpn.com.ovpn", country: "Ireland", region: "Europe", city: "Dublin"},
				"gw1.ewr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ewr1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Newark"},
				"gw1.fra1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.fra1.slickvpn.com.ovpn", country: "Germany", region: "Europe", city: "Frankfurt"},
				"gw1.fra2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.fra2.slickvpn.com.ovpn", country: "Germany", region: "Europe", city: "Frankfurt"},
				"gw1.gru2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.gru2.slickvpn.com.ovpn", country: "Brazil", region: "South America", city: "Sao Paulo"},
				"gw1.grz1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.grz1.slickvpn.com.ovpn", country: "Austria", region: "Europe", city: "Graz"},
				"gw1.had2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.had2.slickvpn.com.ovpn", country: "Sweden", region: "Europe", city: "Halmstad"},
				"gw1.hkg2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.hkg2.slickvpn.com.ovpn", country: "Hong Kong", region: "Asia", city: "Hong Kong"},
				"gw1.iad1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.iad1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Washington"},
				"gw1.iev1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.iev1.slickvpn.com.ovpn", country: "Ukraine", region: "Europe", city: "Kiev"},
				"gw1.iom1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.iom1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "Isle Of Man"},
				"gw1.kiv1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.kiv1.slickvpn.com.ovpn", country: "Moldova", region: "Europe", city: "Chisinau"},
				"gw1.kna1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.kna1.slickvpn.com.ovpn", country: "Chile", region: "South America", city: "Vina del Mar"},
				"gw1.kul1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.kul1.slickvpn.com.ovpn", country: "Malaysia", region: "Asia", city: "Kuala Lumpur"},
				"gw1.las1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.las1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Las Vegas"},
				"gw1.lax1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lax1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Los Angeles"},
				"gw1.lax2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lax2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Los Angeles"},
				"gw1.led1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.led1.slickvpn.com.ovpn", country: "Russian Federation", region: "Europe", city: "St Petersburg"},
				"gw1.lga1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lga1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "New York"},
				"gw1.lga2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lga2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "New York"},
				"gw1.lhr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lhr1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "London"},
				"gw1.lhr2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lhr2.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "London"},
				"gw1.lil1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lil1.slickvpn.com.ovpn", country: "France", region: "Europe", city: "Lille"},
				"gw1.lju1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.lju1.slickvpn.com.ovpn", country: "Slovenia", region: "Europe", city: "Ljubljana"},
				"gw1.mad1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.mad1.slickvpn.com.ovpn", country: "Spain", region: "Europe", city: "Madrid"},
				"gw1.man2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.man2.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "Manchester"},
				"gw1.mci2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.mci2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Kansas City"},
				"gw1.mrn1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.mrn1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Morganton"},
				"gw1.mxp1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.mxp1.slickvpn.com.ovpn", country: "Italy", region: "Europe", city: "Milan"},
				"gw1.mxp2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.mxp2.slickvpn.com.ovpn", country: "Italy", region: "Europe", city: "Milan"},
				"gw1.nrt1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.nrt1.slickvpn.com.ovpn", country: "Japan", region: "Asia", city: "Tokyo"},
				"gw1.nue1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.nue1.slickvpn.com.ovpn", country: "Germany", region: "Europe", city: "Nürnberg"},
				"gw1.ord3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ord3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Chicago"},
				"gw1.ord4.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ord4.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Chicago"},
				"gw1.ost2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.ost2.slickvpn.com.ovpn", country: "Belgium", region: "Europe", city: "Ostend"},
				"gw1.pao1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.pao1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Palo Alto"},
				"gw1.phx2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.phx2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Phoenix"},
				"gw1.prg1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.prg1.slickvpn.com.ovpn", country: "Czech Republic", region: "Europe", city: "Prague"},
				"gw1.prg2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.prg2.slickvpn.com.ovpn", country: "Czech Republic", region: "Europe", city: "Prague"},
				"gw1.rcs1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.rcs1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "Rochester"},
				"gw1.rkv1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.rkv1.slickvpn.com.ovpn", country: "Iceland", region: "Europe", city: "Reykjavik"},
				"gw1.san1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.san1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "San Diego"},
				"gw1.sea1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.sea1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Seattle"},
				"gw1.sea2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.sea2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Seattle"},
				"gw1.sin1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.sin1.slickvpn.com.ovpn", country: "Singapore", region: "Asia", city: "Singapore"},
				"gw1.sin2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.sin2.slickvpn.com.ovpn", country: "Singapore", region: "Asia", city: "Singapore"},
				"gw1.sjc2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.sjc2.slickvpn.com.ovpn", country: "United States", region: "North America", city: "San Jose"},
				"gw1.skg1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.skg1.slickvpn.com.ovpn", country: "Greece", region: "Europe", city: "Thessaloniki"},
				"gw1.sou1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.sou1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "Eastleigh near Southampton"},
				"gw1.stl1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.stl1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "St Louis"},
				"gw1.svo1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.svo1.slickvpn.com.ovpn", country: "Russian Federation", region: "Europe", city: "Moscow"},
				"gw1.svo2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.svo2.slickvpn.com.ovpn", country: "Russian Federation", region: "Europe", city: "Moscow"},
				"gw1.syd1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.syd1.slickvpn.com.ovpn", country: "Australia", region: "Oceania", city: "Sydney"},
				"gw1.syd2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.syd2.slickvpn.com.ovpn", country: "Australia", region: "Oceania", city: "Sydney"},
				"gw1.tll1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.tll1.slickvpn.com.ovpn", country: "Estonia", region: "Europe", city: "Tallinn"},
				"gw1.tlv2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.tlv2.slickvpn.com.ovpn", country: "Israel", region: "Asia", city: "Tel Aviv Yafo"},
				"gw1.tpa1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.tpa1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Tampa"},
				"gw1.trf1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.trf1.slickvpn.com.ovpn", country: "Norway", region: "Europe", city: "Torp"},
				"gw1.waw1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.waw1.slickvpn.com.ovpn", country: "Poland", region: "Europe", city: "Warsaw"},
				"gw1.yei1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.yei1.slickvpn.com.ovpn", country: "Turkey", region: "Asia", city: "Bursa"},
				"gw1.yul1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.yul1.slickvpn.com.ovpn", country: "Canada", region: "North America", city: "Montreal"},
				"gw1.yul2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.yul2.slickvpn.com.ovpn", country: "Canada", region: "North America", city: "Montreal"},
				"gw1.yvr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.yvr1.slickvpn.com.ovpn", country: "Canada", region: "North America", city: "Vancouver"},
				"gw1.yyz1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.yyz1.slickvpn.com.ovpn", country: "Canada", region: "North America", city: "Toronto"},
				"gw1.zrh1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw1.zrh1.slickvpn.com.ovpn", country: "Switzerland", region: "Europe", city: "Zurich"},
				"gw2.ams3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.ams3.slickvpn.com.ovpn", country: "Netherlands", region: "Europe", city: "Amsterdam"},
				"gw2.atl3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.atl3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Atlanta"},
				"gw2.bcn2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.bcn2.slickvpn.com.ovpn", country: "Spain", region: "Europe", city: "Barcelona"},
				"gw2.clt1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.clt1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Charlotte"},
				"gw2.dfw3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.dfw3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Dallas"},
				"gw2.ewr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.ewr1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Newark"},
				"gw2.fra1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.fra1.slickvpn.com.ovpn", country: "Germany", region: "Europe", city: "Frankfurt"},
				"gw2.hou1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.hou1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Houston"},
				"gw2.iad1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.iad1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Washington"},
				"gw2.lhr2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.lhr2.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "London"},
				"gw2.mel1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.mel1.slickvpn.com.ovpn", country: "Australia", region: "Oceania", city: "Melbourne"},
				"gw2.mia3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.mia3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Miami"},
				"gw2.mia4.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.mia4.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Miami"},
				"gw2.mxp2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.mxp2.slickvpn.com.ovpn", country: "Italy", region: "Europe", city: "Milan"},
				"gw2.ord1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.ord1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Chicago"},
				"gw2.ost2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.ost2.slickvpn.com.ovpn", country: "Belgium", region: "Europe", city: "Ostend"},
				"gw2.pao1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.pao1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Palo Alto"},
				"gw2.prg1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.prg1.slickvpn.com.ovpn", country: "Czech Republic", region: "Europe", city: "Prague"},
				"gw2.pty1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.pty1.slickvpn.com.ovpn", country: "Panama", region: "North America", city: "Panama City"},
				"gw2.sin2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.sin2.slickvpn.com.ovpn", country: "Singapore", region: "Asia", city: "Singapore"},
				"gw2.slc1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.slc1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Salt Lake City"},
				"gw2.syd2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.syd2.slickvpn.com.ovpn", country: "Australia", region: "Oceania", city: "Sydney"},
				"gw2.tpe1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.tpe1.slickvpn.com.ovpn", country: "Taiwan", region: "Asia", city: "Taipei"},
				"gw2.yul2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.yul2.slickvpn.com.ovpn", country: "Canada", region: "North America", city: "Montreal"},
				"gw2.yyz1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw2.yyz1.slickvpn.com.ovpn", country: "Canada", region: "North America", city: "Toronto"},
				"gw3.dfw3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.dfw3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Dallas"},
				"gw3.ewr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.ewr1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Newark"},
				"gw3.iad1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.iad1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Washington"},
				"gw3.lax3.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.lax3.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Los Angeles"},
				"gw3.lhr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.lhr1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "London"},
				"gw3.lhr2.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.lhr2.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "London"},
				"gw3.pao1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.pao1.slickvpn.com.ovpn", country: "United States", region: "North America", city: "Palo Alto"},
				"gw3.per1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.per1.slickvpn.com.ovpn", country: "Australia", region: "Oceania", city: "Perth"},
				"gw3.sou1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw3.sou1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "Eastleigh near Southampton"},
				"gw4.lhr1.slickvpn.com": {ovpnURL: "https://www.slickvpn.com/wp-content/themes/slickvpn-theme/vpn-configs/ovpn/gw4.lhr1.slickvpn.com.ovpn", country: "United Kingdom", region: "Europe", city: "London"},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			hostToData, err := parseHTML(testCase.rootNode)

			assert.Equal(t, testCase.hostToData, hostToData)
			assert.ErrorIs(t, err, testCase.errWrapped)
			if testCase.errWrapped != nil {
				assert.EqualError(t, err, testCase.errMessage)
			}
		})
	}
}
