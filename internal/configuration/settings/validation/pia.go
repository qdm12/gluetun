package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func PIAGeoChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func PIAHostnameChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

func PIANameChoices(servers []models.Server) (choices []string) { // TODO remove in v4
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].ServerName
	}
	return makeUnique(choices)
}
