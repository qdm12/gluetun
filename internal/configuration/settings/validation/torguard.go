package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func TorguardCountryChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func TorguardCityChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func TorguardHostnameChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
