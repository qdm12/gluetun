package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func MullvadCountryChoices(servers []models.MullvadServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func MullvadCityChoices(servers []models.MullvadServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func MullvadHostnameChoices(servers []models.MullvadServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

func MullvadISPChoices(servers []models.MullvadServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].ISP
	}
	return makeUnique(choices)
}
