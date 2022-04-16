package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func PerfectprivacyCityChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}
