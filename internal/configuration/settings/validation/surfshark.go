package validation

import (
	"sort"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

func SurfsharkRegionChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Region
	}
	return makeUnique(choices)
}

func SurfsharkCountryChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Country
	}
	return makeUnique(choices)
}

func SurfsharkCityChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].City
	}
	return makeUnique(choices)
}

func SurfsharkHostnameChoices(servers []models.Server) (choices []string) {
	choices = make([]string, len(servers))
	for i := range servers {
		choices[i] = servers[i].Hostname
	}
	return makeUnique(choices)
}

// TODO remove in v4.
func SurfsharkRetroLocChoices() (choices []string) {
	locationData := constants.SurfsharkLocationData()
	choices = make([]string, 0, len(locationData))
	seen := make(map[string]struct{}, len(locationData))
	for _, data := range locationData {
		if _, ok := seen[data.RetroLoc]; ok {
			continue
		}
		seen[data.RetroLoc] = struct{}{}
		choices = append(choices, data.RetroLoc)
	}

	sort.Slice(choices, func(i, j int) bool {
		return choices[i] < choices[j]
	})

	return choices
}
