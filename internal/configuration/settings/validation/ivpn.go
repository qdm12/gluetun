package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func IvpnCountryChoices(servers []models.IvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func IvpnCityChoices(servers []models.IvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func IvpnISPChoices(servers []models.IvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].ISP
	}
	return makeUnique(choices)
}

func IvpnHostnameChoices(servers []models.IvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
