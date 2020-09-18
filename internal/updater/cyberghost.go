package updater

import (
	"context"
	"fmt"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateCyberghost(ctx context.Context) (err error) {
	servers, err := findCyberghostServers(ctx, u.lookupIP)
	if err != nil {
		return err
	}
	if u.options.Stdout {
		u.println(stringifyCyberghostServers(servers))
	}
	u.servers.Cyberghost.Timestamp = u.timeNow().Unix()
	u.servers.Cyberghost.Servers = servers
	return nil
}

func findCyberghostServers(ctx context.Context, lookupIP lookupIPFunc) (servers []models.CyberghostServer, err error) {
	groups := getCyberghostGroups()
	allCountryCodes := getCountryCodes()
	cyberghostCountryCodes := getCyberghostSubdomainToRegion()
	possibleCountryCodes := mergeCountryCodes(cyberghostCountryCodes, allCountryCodes)

	results := make(chan models.CyberghostServer)
	const maxGoroutines = 10
	guard := make(chan struct{}, maxGoroutines)
	for groupID, groupName := range groups {
		for countryCode, region := range possibleCountryCodes {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			const domain = "cg-dialup.net"
			host := fmt.Sprintf("%s-%s.%s", groupID, countryCode, domain)
			guard <- struct{}{}
			go tryCyberghostHostname(ctx, lookupIP, host, groupName, region, results)
			<-guard
		}
	}
	for i := 0; i < len(groups)*len(possibleCountryCodes); i++ {
		server := <-results
		if server.IPs == nil {
			continue
		}
		servers = append(servers, server)
	}
	if err := ctx.Err(); err != nil {
		return servers, err
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	return servers, nil
}

func tryCyberghostHostname(ctx context.Context, lookupIP lookupIPFunc,
	host, groupName, region string,
	results chan<- models.CyberghostServer) {
	IPs, err := resolveRepeat(ctx, lookupIP, host, 2)
	if err != nil || len(IPs) == 0 {
		results <- models.CyberghostServer{}
		return
	}
	results <- models.CyberghostServer{
		Region: region,
		Group:  groupName,
		IPs:    IPs,
	}
}

//nolint:goconst
func stringifyCyberghostServers(servers []models.CyberghostServer) (s string) {
	s = "func CyberghostServers() []models.CyberghostServer {\n"
	s += "	return []models.CyberghostServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}

func getCyberghostGroups() map[string]string {
	return map[string]string{
		"87-1": "Premium UDP Europe",
		"94-1": "Premium UDP USA",
		"95-1": "Premium UDP Asia",
		"87-8": "NoSpy UDP Europe",
		"97-1": "Premium TCP Europe",
		"93-1": "Premium TCP USA",
		"96-1": "Premium TCP Asia",
		"97-8": "NoSpy TCP Europe",
	}
}

func getCyberghostSubdomainToRegion() map[string]string { //nolint:dupl
	return map[string]string{
		"af": "Afghanistan",
		"ax": "Aland Islands",
		"al": "Albania",
		"dz": "Algeria",
		"as": "American Samoa",
		"ad": "Andorra",
		"ao": "Angola",
		"ai": "Anguilla",
		"aq": "Antarctica",
		"ag": "Antigua and Barbuda",
		"ar": "Argentina",
		"am": "Armenia",
		"aw": "Aruba",
		"au": "Australia",
		"at": "Austria",
		"az": "Azerbaijan",
		"bs": "Bahamas",
		"bh": "Bahrain",
		"bd": "Bangladesh",
		"bb": "Barbados",
		"by": "Belarus",
		"be": "Belgium",
		"bz": "Belize",
		"bj": "Benin",
		"bm": "Bermuda",
		"bt": "Bhutan",
		"bo": "Bolivia",
		"bq": "Bonaire",
		"ba": "Bosnia and Herzegovina",
		"bw": "Botswana",
		"bv": "Bouvet Island",
		"br": "Brazil",
		"io": "British Indian Ocean Territory",
		"vg": "British Virgin Islands",
		"bn": "Brunei Darussalam",
		"bg": "Bulgaria",
		"bf": "Burkina Faso",
		"bi": "Burundi",
		"kh": "Cambodia",
		"cm": "Cameroon",
		"ca": "Canada",
		"cv": "Cape Verde",
		"ky": "Cayman Islands",
		"cf": "Central African Republic",
		"td": "Chad",
		"cl": "Chile",
		"cn": "China",
		"cx": "Christmas Island",
		"cc": "Cocos Islands",
		"co": "Colombia",
		"km": "Comoros",
		"cg": "Congo",
		"ck": "Cook Islands",
		"cr": "Costa Rica",
		"ci": "Cote d'Ivoire",
		"hr": "Croatia",
		"cu": "Cuba",
		"cw": "Curacao",
		"cy": "Cyprus",
		"cz": "Czech Republic",
		"cd": "Democratic Republic of the Congo",
		"dk": "Denmark",
		"dj": "Djibouti",
		"dm": "Dominica",
		"do": "Dominican Republic",
		"ec": "Ecuador",
		"eg": "Egypt",
		"sv": "El Salvador",
		"gq": "Equatorial Guinea",
		"er": "Eritrea",
		"ee": "Estonia",
		"et": "Ethiopia",
		"fk": "Falkland Islands",
		"fo": "Faroe Islands",
		"fj": "Fiji",
		"fi": "Finland",
		"fr": "France",
		"gf": "French Guiana",
		"pf": "French Polynesia",
		"tf": "French Southern Territories",
		"ga": "Gabon",
		"gm": "Gambia",
		"ge": "Georgia",
		"de": "Germany",
		"gh": "Ghana",
		"gi": "Gibraltar",
		"gr": "Greece",
		"gl": "Greenland",
		"gd": "Grenada",
		"gp": "Guadeloupe",
		"gu": "Guam",
		"gt": "Guatemala",
		"gg": "Guernsey",
		"gw": "Guinea-Bissau",
		"gn": "Guinea",
		"gy": "Guyana",
		"ht": "Haiti",
		"hm": "Heard Island and McDonald Islands",
		"hn": "Honduras",
		"hk": "Hong Kong",
		"hu": "Hungary",
		"is": "Iceland",
		"in": "India",
		"id": "Indonesia",
		"ir": "Iran",
		"iq": "Iraq",
		"ie": "Ireland",
		"im": "Isle of Man",
		"il": "Israel",
		"it": "Italy",
		"jm": "Jamaica",
		"jp": "Japan",
		"je": "Jersey",
		"jo": "Jordan",
		"kz": "Kazakhstan",
		"ke": "Kenya",
		"ki": "Kiribati",
		"kr": "Korea",
		"kw": "Kuwait",
		"kg": "Kyrgyzstan",
		"la": "Lao People's Democratic Republic",
		"lv": "Latvia",
		"lb": "Lebanon",
		"ls": "Lesotho",
		"lr": "Liberia",
		"ly": "Libya",
		"li": "Liechtenstein",
		"lt": "Lithuania",
		"lu": "Luxembourg",
		"mo": "Macao",
		"mk": "Macedonia",
		"mg": "Madagascar",
		"mw": "Malawi",
		"my": "Malaysia",
		"mv": "Maldives",
		"ml": "Mali",
		"mt": "Malta",
		"mh": "Marshall Islands",
		"mq": "Martinique",
		"mr": "Mauritania",
		"mu": "Mauritius",
		"yt": "Mayotte",
		"mx": "Mexico",
		"fm": "Micronesia",
		"md": "Moldova",
		"mc": "Monaco",
		"mn": "Mongolia",
		"me": "Montenegro",
		"ms": "Montserrat",
		"ma": "Morocco",
		"mz": "Mozambique",
		"mm": "Myanmar",
		"na": "Namibia",
		"nr": "Nauru",
		"np": "Nepal",
		"nl": "Netherlands",
		"nc": "New Caledonia",
		"nz": "New Zealand",
		"ni": "Nicaragua",
		"ne": "Niger",
		"ng": "Nigeria",
		"nu": "Niue",
		"nf": "Norfolk Island",
		"mp": "Northern Mariana Islands",
		"no": "Norway",
		"om": "Oman",
		"pk": "Pakistan",
		"pw": "Palau",
		"ps": "Palestine, State of",
		"pa": "Panama",
		"pg": "Papua New Guinea",
		"py": "Paraguay",
		"pe": "Peru",
		"ph": "Philippines",
		"pn": "Pitcairn",
		"pl": "Poland",
		"pt": "Portugal",
		"pr": "Puerto Rico",
		"qa": "Qatar",
		"re": "Reunion",
		"ro": "Romania",
		"ru": "Russian Federation",
		"rw": "Rwanda",
		"bl": "Saint Barthelemy",
		"sh": "Saint Helena",
		"kn": "Saint Kitts and Nevis",
		"lc": "Saint Lucia",
		"mf": "Saint Martin",
		"pm": "Saint Pierre and Miquelon",
		"vc": "Saint Vincent and the Grenadines",
		"ws": "Samoa",
		"sm": "San Marino",
		"st": "Sao Tome and Principe",
		"sa": "Saudi Arabia",
		"sn": "Senegal",
		"rs": "Serbia",
		"sc": "Seychelles",
		"sl": "Sierra Leone",
		"sg": "Singapore",
		"sx": "Sint Maarten",
		"sk": "Slovakia",
		"si": "Slovenia",
		"sb": "Solomon Islands",
		"so": "Somalia",
		"za": "South Africa",
		"gs": "South Georgia and the South Sandwich Islands",
		"ss": "South Sudan",
		"es": "Spain",
		"lk": "Sri Lanka",
		"sd": "Sudan",
		"sr": "Suriname",
		"sj": "Svalbard and Jan Mayen",
		"sz": "Swaziland",
		"se": "Sweden",
		"ch": "Switzerland",
		"sy": "Syrian Arab Republic",
		"tw": "Taiwan",
		"tj": "Tajikistan",
		"tz": "Tanzania",
		"th": "Thailand",
		"tl": "Timor-Leste",
		"tg": "Togo",
		"tk": "Tokelau",
		"to": "Tonga",
		"tt": "Trinidad and Tobago",
		"tn": "Tunisia",
		"tr": "Turkey",
		"tm": "Turkmenistan",
		"tc": "Turks and Caicos Islands",
		"tv": "Tuvalu",
		"ug": "Uganda",
		"ua": "Ukraine",
		"ae": "United Arab Emirates",
		"gb": "United Kingdom",
		"um": "United States Minor Outlying Islands",
		"us": "United States",
		"uy": "Uruguay",
		"vi": "US Virgin Islands",
		"uz": "Uzbekistan",
		"vu": "Vanuatu",
		"va": "Vatican City State",
		"ve": "Venezuela",
		"vn": "Vietnam",
		"wf": "Wallis and Futuna",
		"eh": "Western Sahara",
		"ye": "Yemen",
		"zm": "Zambia",
		"zw": "Zimbabwe",
	}
}
