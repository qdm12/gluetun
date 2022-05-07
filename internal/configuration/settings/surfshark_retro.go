package settings

import (
	"strings"

	"github.com/qdm12/gluetun/internal/provider/surfshark/servers"
)

func surfsharkRetroRegion(selection ServerSelection) (
	updatedSelection ServerSelection) {
	locationData := servers.LocationData()

	retroToLocation := make(map[string]servers.ServerLocation, len(locationData))
	for _, data := range locationData {
		if data.RetroLoc == "" {
			continue
		}
		retroToLocation[strings.ToLower(data.RetroLoc)] = data
	}

	for i, region := range selection.Regions {
		location, ok := retroToLocation[region]
		if !ok {
			continue
		}
		selection.Regions[i] = strings.ToLower(location.Region)
		selection.Countries = append(selection.Countries, strings.ToLower(location.Country))
		selection.Cities = append(selection.Cities, strings.ToLower(location.City)) // even empty string
		selection.Hostnames = append(selection.Hostnames, location.Hostname)
	}

	selection.Regions = dedupSlice(selection.Regions)
	selection.Countries = dedupSlice(selection.Countries)
	selection.Cities = dedupSlice(selection.Cities)
	selection.Hostnames = dedupSlice(selection.Hostnames)

	return selection
}

func dedupSlice(slice []string) (deduped []string) {
	if slice == nil {
		return nil
	}

	deduped = make([]string, 0, len(slice))
	seen := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			deduped = append(deduped, s)
		}
	}

	return deduped
}
