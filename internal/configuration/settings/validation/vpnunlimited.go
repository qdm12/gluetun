package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func VPNUnlimitedCountryChoices(servers []models.VPNUnlimitedServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func VPNUnlimitedCityChoices(servers []models.VPNUnlimitedServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func VPNUnlimitedHostnameChoices(servers []models.VPNUnlimitedServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
