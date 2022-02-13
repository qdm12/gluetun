package validation

import (
	"github.com/qdm12/gluetun/internal/models"
)

func VyprvpnRegionChoices(servers []models.VyprvpnServer) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}
