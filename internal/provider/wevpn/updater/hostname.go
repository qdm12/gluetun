package updater

import "strings"

func getHostnameFromCity(city string) (hostname string) {
	host := strings.ToLower(city)
	host = strings.ReplaceAll(host, ".", "")
	host = strings.ReplaceAll(host, " ", "")

	specialCases := map[string]string{
		"washingtondc": "washington",
		"mexicocity":   "mexico",
		"denizli":      "bursa",
		"sibu":         "kualalumpur",
		"kiev":         "kyiv",
		"stpetersburg": "petersburg",
	}
	if specialHost, ok := specialCases[host]; ok {
		host = specialHost
	}

	hostname = host + ".wevpn.com"
	return hostname
}
