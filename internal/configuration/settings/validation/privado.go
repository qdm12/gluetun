package validation

import "github.com/qdm12/gluetun/internal/models"

func PrivadoCountryChoices(servers []models.PrivadoServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func PrivadoRegionChoices(servers []models.PrivadoServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func PrivadoCityChoices(servers []models.PrivadoServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func PrivadoHostnameChoices(servers []models.PrivadoServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}
