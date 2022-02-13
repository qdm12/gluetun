package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func PIAGeoChoices(servers []models.PIAServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func PIAHostnameChoices(servers []models.PIAServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

func PIANameChoices(servers []models.PIAServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].ServerName
	}
	return makeUnique(choices)
}
