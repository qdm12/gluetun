package validation

import (
	"github.com/qdm12/gluetun/internal/provider/surfshark/servers"
)

// TODO remove in v4.
func SurfsharkRetroLocChoices() (choices []string) {
	locationData := servers.LocationData()
	choices = make([]string, 0, len(locationData))
	seen := make(map[string]struct{}, len(locationData))
	for _, data := range locationData {
		if _, ok := seen[data.RetroLoc]; ok {
			continue
		}
		seen[data.RetroLoc] = struct{}{}
		choices = sortedInsert(choices, data.RetroLoc)
	}

	return choices
}
