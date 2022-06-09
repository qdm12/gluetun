package updater

import (
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/constants/vpn"
	"github.com/qdm12/gluetun/internal/models"
)

func getHostToServer() (hts hostToServer, warnings []string) {
	shortHTS := map[string]models.Server{
		"ae": {},
		"ar": {},
		"at": {},
		"au-syd": {
			City: "Sydney",
		},
		"ba": {},
		"be": {},
		"bg": {},
		"br": {},
		"by": {},
		"ca-tr": {
			City: "Toronto",
		},
		"ca-vn": {
			City: "Vancouver",
		},
		"ca": {},
		"ch": {},
		"cr": {},
		"cy": {},
		"cz": {},
		"de-dus": {
			City: "Düsseldorf",
		},
		"de": {},
		"dk": {},
		"ee": {},
		"es": {},
		"fi": {},
		"fr-rbx": {
			City: "Roubaix",
		},
		"fr": {},
		"gr": {},
		"hr": {},
		"hu": {},
		"ie-dub": {
			City: "Dublin",
		},
		"il": {},
		"im": {},
		"in-ka": {
			City: "Karnataka",
		},
		"in": {},
		"is": {},
		"it-mil": {
			City: "Milan",
		},
		"jp":  {},
		"kr":  {},
		"lt":  {},
		"lu":  {},
		"lv":  {},
		"ly":  {},
		"md":  {},
		"mx":  {},
		"mys": {},
		"nl":  {},
		"no":  {},
		"nz":  {},
		"om":  {},
		"pl":  {},
		"pt":  {},
		"ro":  {},
		"rs":  {},
		"se":  {},
		"sg-free": {
			Free: true,
		},
		"sg": {},
		"si": {},
		"sk": {},
		"th": {},
		"tr": {},
		"uk-cv": {
			City: "London",
		},
		"uk-lon": {
			City: "London",
		},
		"uk": {},
		"us-chi": {
			City: "Chicago",
		},
		"us-dal": {
			City: "Dallas",
		},
		"us-den": {
			City: "Denver",
		},
		"us-hou": {
			City: "Houston",
		},
		"us-la": {
			City: "Los Angeles",
		},
		"us-lv": {
			City: "Las Vegas",
		},
		"us-mia": {
			City: "Miami",
		},
		"us-ny-free": {
			City: "New York",
			Free: true,
		},
		"us-ny": {
			City: "New York",
		},
		"us-sea": {
			City: "Seattle",
		},
		"us-sf": {
			City: "San Francisco",
		},
		"us-sl": {
			City: "Saint Louis",
		},
		"us-slc": {
			City: "Salt Lake City",
		},
		"us-stream": {
			Stream: true,
		},
		"us": {},
		"vn": {},
		"za": {},
	}

	hts = make(hostToServer, len(shortHTS))

	countryCodesMap := constants.CountryCodes()
	for shortHost, server := range shortHTS {
		server.VPN = vpn.OpenVPN
		server.UDP = true
		server.Hostname = shortHost + ".vpnunlimitedapp.com"
		countryCode := strings.Split(shortHost, "-")[0]
		country, ok := countryCodesMap[countryCode]
		if !ok {
			warnings = append(warnings, "country code not found: "+countryCode)
			continue
		}
		server.Country = country
		hts[server.Hostname] = server
	}

	return hts, warnings
}
