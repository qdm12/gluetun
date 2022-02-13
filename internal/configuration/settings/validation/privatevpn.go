package validation

import "github.com/qdm12/gluetun/internal/models"

func PrivatevpnCountryChoices(servers []models.PrivatevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func PrivatevpnCityChoices(servers []models.PrivatevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func PrivatevpnHostnameChoices(servers []models.PrivatevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
