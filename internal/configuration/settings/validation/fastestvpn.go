package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func FastestvpnCountriesChoices(servers []models.FastestvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func FastestvpnHostnameChoices(servers []models.FastestvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
