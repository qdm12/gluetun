package validation

import "github.com/qdm12/gluetun/internal/models"

func PurevpnRegionChoices(servers []models.PurevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func PurevpnCountryChoices(servers []models.PurevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func PurevpnCityChoices(servers []models.PurevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func PurevpnHostnameChoices(servers []models.PurevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
