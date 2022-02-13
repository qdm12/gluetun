package validation

import "github.com/qdm12/gluetun/internal/models"

func ProtonvpnCountryChoices(servers []models.ProtonvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func ProtonvpnRegionChoices(servers []models.ProtonvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func ProtonvpnCityChoices(servers []models.ProtonvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func ProtonvpnNameChoices(servers []models.ProtonvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Name
	}
	return makeUnique(choices)
}

func ProtonvpnHostnameChoices(servers []models.ProtonvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
