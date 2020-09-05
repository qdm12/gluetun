package updater

import (
	"context"
	"sort"

	"github.com/qdm12/gluetun/internal/models"
)

func (u *updater) updateWindscribe(ctx context.Context) {
	servers := findWindscribeServers(ctx, u.lookupIP)
	if u.options.Stdout {
		u.println(stringifyWindscribeServers(servers))
	}
	u.servers.Windscribe.Timestamp = u.timeNow().Unix()
	u.servers.Windscribe.Servers = servers
}

func findWindscribeServers(ctx context.Context, lookupIP lookupIPFunc) (servers []models.WindscribeServer) {
	allCountryCodes := getCountryCodes()
	windscribeCountryCodes := getWindscribeSubdomainToRegion()
	possibleCountryCodes := mergeCountryCodes(windscribeCountryCodes, allCountryCodes)
	const domain = "windscribe.com"
	for countryCode, region := range possibleCountryCodes {
		host := countryCode + "." + domain
		ips, err := resolveRepeat(ctx, lookupIP, host, 2)
		if err != nil || len(ips) == 0 {
			continue
		}
		servers = append(servers, models.WindscribeServer{
			Region: region,
			IPs:    ips,
		})
	}
	sort.Slice(servers, func(i, j int) bool {
		return servers[i].Region < servers[j].Region
	})
	return servers
}

func mergeCountryCodes(base, extend map[string]string) (merged map[string]string) {
	merged = make(map[string]string, len(base))
	for countryCode, region := range base {
		merged[countryCode] = region
	}
	for countryCode := range base {
		delete(extend, countryCode)
	}
	for countryCode, region := range extend {
		merged[countryCode] = region
	}
	return merged
}

func stringifyWindscribeServers(servers []models.WindscribeServer) (s string) {
	s = "func WindscribeServers() []models.WindscribeServer {\n"
	s += "	return []models.WindscribeServer{\n"
	for _, server := range servers {
		s += "		" + server.String() + ",\n"
	}
	s += "	}\n"
	s += "}"
	return s
}

func getWindscribeSubdomainToRegion() map[string]string {
	return map[string]string{
		"al":         "Albania",
		"ar":         "Argentina",
		"au":         "Australia",
		"at":         "Austria",
		"az":         "Azerbaijan",
		"be":         "Belgium",
		"ba":         "Bosnia",
		"br":         "Brazil",
		"bg":         "Bulgaria",
		"ca":         "Canada East",
		"ca-west":    "Canada West",
		"co":         "Colombia",
		"hr":         "Croatia",
		"cy":         "Cyprus",
		"cz":         "Czech republic",
		"dk":         "Denmark",
		"ee":         "Estonia",
		"aq":         "Fake antarctica",
		"fi":         "Finland",
		"fr":         "France",
		"ge":         "Georgia",
		"de":         "Germany",
		"gr":         "Greece",
		"hk":         "Hong kong",
		"hu":         "Hungary",
		"is":         "Iceland",
		"in":         "India",
		"id":         "Indonesia",
		"ie":         "Ireland",
		"il":         "Israel",
		"it":         "Italy",
		"jp":         "Japan",
		"lv":         "Latvia",
		"lt":         "Lithuania",
		"mk":         "Macedonia",
		"my":         "Malaysia",
		"mx":         "Mexico",
		"md":         "Moldova",
		"nl":         "Netherlands",
		"nz":         "New zealand",
		"no":         "Norway",
		"ph":         "Philippines",
		"pl":         "Poland",
		"pt":         "Portugal",
		"ro":         "Romania",
		"ru":         "Russia",
		"rs":         "Serbia",
		"sg":         "Singapore",
		"sk":         "Slovakia",
		"si":         "Slovenia",
		"za":         "South Africa",
		"kr":         "South Korea",
		"es":         "Spain",
		"se":         "Sweden",
		"ch":         "Switzerland",
		"th":         "Thailand",
		"tn":         "Tunisia",
		"tr":         "Turkey",
		"ua":         "Ukraine",
		"ae":         "United Arab Emirates",
		"uk":         "United Kingdom",
		"us-central": "US Central",
		"us-east":    "US East",
		"us-west":    "US West",
		"vn":         "Vietnam",
		"wf-ca":      "Windflix CA",
		"wf-jp":      "Windflix JP",
		"wf-uk":      "Windflix UK",
		"wf-us":      "Windflix US",
	}
}
