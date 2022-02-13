package validation

import "github.com/qdm12/gluetun/internal/models"

func WevpnCityChoices(servers []models.WevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func WevpnHostnameChoices(servers []models.WevpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
