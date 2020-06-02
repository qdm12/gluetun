package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
)

func main() {
	ctx := context.Background()
	os.Exit(_main(ctx))
}

func _main(ctx context.Context) int {
	fmt.Println("Host finder for Cyberghost")
	resolverAddress := flag.String("resolver", "1.1.1.1", "DNS Resolver IP address to use")
	flag.Parse()

	resolver := newResolver(*resolverAddress)
	lookupIP := newLookupIP(resolver)

	const domain = "cg-dialup.net"
	groups := getCyberghostGroups()
	countryCodes := getCountryCodes()
	type result struct {
		groupName string
		region    string
		subdomain string
		exists    bool
	}
	resultsChannel := make(chan result)
	const maxGoroutines = 10
	guard := make(chan struct{}, maxGoroutines)
	fmt.Print("Subdomains found: ")
	for groupName, groupID := range groups {
		for country, countryCode := range countryCodes {
			go func(groupName, groupID, country, countryCode string) {
				r := result{
					region:    country,
					groupName: groupName,
					subdomain: fmt.Sprintf("%s-%s", groupID, countryCode),
				}
				fqdn := fmt.Sprintf("%s.%s", r.subdomain, domain)
				guard <- struct{}{}
				ips, err := lookupIP(ctx, fqdn)
				<-guard
				if err == nil && len(ips) > 0 {
					r.exists = true
				}
				resultsChannel <- r
			}(groupName, groupID, country, countryCode)
		}
	}
	results := make([]result, len(groups)*len(countryCodes))
	for i := range results {
		results[i] = <-resultsChannel
		fmt.Printf("%s ", results[i].subdomain)
	}
	fmt.Print("\n\n")
	sort.Slice(results, func(i, j int) bool {
		return results[i].region < results[j].region
	})
	for _, r := range results {
		if r.exists {
			// Use in resolver program
			fmt.Printf("{subdomain: %q, region: %q, group: %q},\n", r.subdomain, r.region, r.groupName)
		}
	}
	return 0
}

func newResolver(ip string) *net.Resolver {
	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "udp", net.JoinHostPort(ip, "53"))
		},
	}
}

type lookupIPFunc func(ctx context.Context, host string) (ips []net.IP, err error)

func newLookupIP(r *net.Resolver) lookupIPFunc {
	return func(ctx context.Context, host string) (ips []net.IP, err error) {
		addresses, err := r.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err
		}
		ips = make([]net.IP, len(addresses))
		for i := range addresses {
			ips[i] = addresses[i].IP
		}
		return ips, nil
	}
}

func getCyberghostGroups() map[string]string {
	return map[string]string{
		"Premium UDP Europe": "87-1",
		"Premium UDP USA":    "94-1",
		"Premium UDP Asia":   "95-1",
		"NoSpy UDP Europe":   "87-8",
		"Premium TCP Europe": "97-1",
		"Premium TCP USA":    "93-1",
		"Premium TCP Asia":   "96-1",
		"NoSpy TCP Europe":   "97-8",
	}
}

func getCountryCodes() map[string]string {
	return map[string]string{
		"Afghanistan":                       "af",
		"Aland Islands":                     "ax",
		"Albania":                           "al",
		"Algeria":                           "dz",
		"American Samoa":                    "as",
		"Andorra":                           "ad",
		"Angola":                            "ao",
		"Anguilla":                          "ai",
		"Antarctica":                        "aq",
		"Antigua and Barbuda":               "ag",
		"Argentina":                         "ar",
		"Armenia":                           "am",
		"Aruba":                             "aw",
		"Australia":                         "au",
		"Austria":                           "at",
		"Azerbaijan":                        "az",
		"Bahamas":                           "bs",
		"Bahrain":                           "bh",
		"Bangladesh":                        "bd",
		"Barbados":                          "bb",
		"Belarus":                           "by",
		"Belgium":                           "be",
		"Belize":                            "bz",
		"Benin":                             "bj",
		"Bermuda":                           "bm",
		"Bhutan":                            "bt",
		"Bolivia":                           "bo",
		"Bonaire":                           "bq",
		"Bosnia and Herzegovina":            "ba",
		"Botswana":                          "bw",
		"Bouvet Island":                     "bv",
		"Brazil":                            "br",
		"British Indian Ocean Territory":    "io",
		"British Virgin Islands":            "vg",
		"Brunei Darussalam":                 "bn",
		"Bulgaria":                          "bg",
		"Burkina Faso":                      "bf",
		"Burundi":                           "bi",
		"Cambodia":                          "kh",
		"Cameroon":                          "cm",
		"Canada":                            "ca",
		"Cape Verde":                        "cv",
		"Cayman Islands":                    "ky",
		"Central African Republic":          "cf",
		"Chad":                              "td",
		"Chile":                             "cl",
		"China":                             "cn",
		"Christmas Island":                  "cx",
		"Cocos Islands":                     "cc",
		"Colombia":                          "co",
		"Comoros":                           "km",
		"Congo":                             "cg",
		"Cook Islands":                      "ck",
		"Costa Rica":                        "cr",
		"Cote d'Ivoire":                     "ci",
		"Croatia":                           "hr",
		"Cuba":                              "cu",
		"Curacao":                           "cw",
		"Cyprus":                            "cy",
		"Czech Republic":                    "cz",
		"Democratic Republic of the Congo":  "cd",
		"Denmark":                           "dk",
		"Djibouti":                          "dj",
		"Dominica":                          "dm",
		"Dominican Republic":                "do",
		"Ecuador":                           "ec",
		"Egypt":                             "eg",
		"El Salvador":                       "sv",
		"Equatorial Guinea":                 "gq",
		"Eritrea":                           "er",
		"Estonia":                           "ee",
		"Ethiopia":                          "et",
		"Falkland Islands":                  "fk",
		"Faroe Islands":                     "fo",
		"Fiji":                              "fj",
		"Finland":                           "fi",
		"France":                            "fr",
		"French Guiana":                     "gf",
		"French Polynesia":                  "pf",
		"French Southern Territories":       "tf",
		"Gabon":                             "ga",
		"Gambia":                            "gm",
		"Georgia":                           "ge",
		"Germany":                           "de",
		"Ghana":                             "gh",
		"Gibraltar":                         "gi",
		"Greece":                            "gr",
		"Greenland":                         "gl",
		"Grenada":                           "gd",
		"Guadeloupe":                        "gp",
		"Guam":                              "gu",
		"Guatemala":                         "gt",
		"Guernsey":                          "gg",
		"Guinea-Bissau":                     "gw",
		"Guinea":                            "gn",
		"Guyana":                            "gy",
		"Haiti":                             "ht",
		"Heard Island and McDonald Islands": "hm",
		"Honduras":                          "hn",
		"Hong Kong":                         "hk",
		"Hungary":                           "hu",
		"Iceland":                           "is",
		"India":                             "in",
		"Indonesia":                         "id",
		"Iran":                              "ir",
		"Iraq":                              "iq",
		"Ireland":                           "ie",
		"Isle of Man":                       "im",
		"Israel":                            "il",
		"Italy":                             "it",
		"Jamaica":                           "jm",
		"Japan":                             "jp",
		"Jersey":                            "je",
		"Jordan":                            "jo",
		"Kazakhstan":                        "kz",
		"Kenya":                             "ke",
		"Kiribati":                          "ki",
		"Korea":                             "kr",
		"Kuwait":                            "kw",
		"Kyrgyzstan":                        "kg",
		"Lao People's Democratic Republic":  "la",
		"Latvia":                            "lv",
		"Lebanon":                           "lb",
		"Lesotho":                           "ls",
		"Liberia":                           "lr",
		"Libya":                             "ly",
		"Liechtenstein":                     "li",
		"Lithuania":                         "lt",
		"Luxembourg":                        "lu",
		"Macao":                             "mo",
		"Macedonia":                         "mk",
		"Madagascar":                        "mg",
		"Malawi":                            "mw",
		"Malaysia":                          "my",
		"Maldives":                          "mv",
		"Mali":                              "ml",
		"Malta":                             "mt",
		"Marshall Islands":                  "mh",
		"Martinique":                        "mq",
		"Mauritania":                        "mr",
		"Mauritius":                         "mu",
		"Mayotte":                           "yt",
		"Mexico":                            "mx",
		"Micronesia":                        "fm",
		"Moldova":                           "md",
		"Monaco":                            "mc",
		"Mongolia":                          "mn",
		"Montenegro":                        "me",
		"Montserrat":                        "ms",
		"Morocco":                           "ma",
		"Mozambique":                        "mz",
		"Myanmar":                           "mm",
		"Namibia":                           "na",
		"Nauru":                             "nr",
		"Nepal":                             "np",
		"Netherlands":                       "nl",
		"New Caledonia":                     "nc",
		"New Zealand":                       "nz",
		"Nicaragua":                         "ni",
		"Niger":                             "ne",
		"Nigeria":                           "ng",
		"Niue":                              "nu",
		"Norfolk Island":                    "nf",
		"Northern Mariana Islands":          "mp",
		"Norway":                            "no",
		"Oman":                              "om",
		"Pakistan":                          "pk",
		"Palau":                             "pw",
		"Palestine, State of":               "ps",
		"Panama":                            "pa",
		"Papua New Guinea":                  "pg",
		"Paraguay":                          "py",
		"Peru":                              "pe",
		"Philippines":                       "ph",
		"Pitcairn":                          "pn",
		"Poland":                            "pl",
		"Portugal":                          "pt",
		"Puerto Rico":                       "pr",
		"Qatar":                             "qa",
		"Reunion":                           "re",
		"Romania":                           "ro",
		"Russian Federation":                "ru",
		"Rwanda":                            "rw",
		"Saint Barthelemy":                  "bl",
		"Saint Helena":                      "sh",
		"Saint Kitts and Nevis":             "kn",
		"Saint Lucia":                       "lc",
		"Saint Martin":                      "mf",
		"Saint Pierre and Miquelon":         "pm",
		"Saint Vincent and the Grenadines":  "vc",
		"Samoa":                             "ws",
		"San Marino":                        "sm",
		"Sao Tome and Principe":             "st",
		"Saudi Arabia":                      "sa",
		"Senegal":                           "sn",
		"Serbia":                            "rs",
		"Seychelles":                        "sc",
		"Sierra Leone":                      "sl",
		"Singapore":                         "sg",
		"Sint Maarten":                      "sx",
		"Slovakia":                          "sk",
		"Slovenia":                          "si",
		"Solomon Islands":                   "sb",
		"Somalia":                           "so",
		"South Africa":                      "za",
		"South Georgia and the South Sandwich Islands": "gs",
		"South Sudan":                          "ss",
		"Spain":                                "es",
		"Sri Lanka":                            "lk",
		"Sudan":                                "sd",
		"Suriname":                             "sr",
		"Svalbard and Jan Mayen":               "sj",
		"Swaziland":                            "sz",
		"Sweden":                               "se",
		"Switzerland":                          "ch",
		"Syrian Arab Republic":                 "sy",
		"Taiwan":                               "tw",
		"Tajikistan":                           "tj",
		"Tanzania":                             "tz",
		"Thailand":                             "th",
		"Timor-Leste":                          "tl",
		"Togo":                                 "tg",
		"Tokelau":                              "tk",
		"Tonga":                                "to",
		"Trinidad and Tobago":                  "tt",
		"Tunisia":                              "tn",
		"Turkey":                               "tr",
		"Turkmenistan":                         "tm",
		"Turks and Caicos Islands":             "tc",
		"Tuvalu":                               "tv",
		"Uganda":                               "ug",
		"Ukraine":                              "ua",
		"United Arab Emirates":                 "ae",
		"United Kingdom":                       "gb",
		"United States Minor Outlying Islands": "um",
		"United States":                        "us",
		"Uruguay":                              "uy",
		"US Virgin Islands":                    "vi",
		"Uzbekistan":                           "uz",
		"Vanuatu":                              "vu",
		"Vatican City State":                   "va",
		"Venezuela":                            "ve",
		"Vietnam":                              "vn",
		"Wallis and Futuna":                    "wf",
		"Western Sahara":                       "eh",
		"Yemen":                                "ye",
		"Zambia":                               "zm",
		"Zimbabwe":                             "zw",
	}
}
