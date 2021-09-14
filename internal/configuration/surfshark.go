package configuration

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
	"github.com/qdm12/golibs/params"
)

func (settings *Provider) readSurfshark(r reader) (err error) {
	settings.Name = constants.Surfshark
	servers := r.servers.GetSurfshark()

	settings.ServerSelection.TargetIP, err = readTargetIP(r.env)
	if err != nil {
		return err
	}

	settings.ServerSelection.Countries, err = r.env.CSVInside("COUNTRY", constants.SurfsharkCountryChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable COUNTRY: %w", err)
	}

	settings.ServerSelection.Cities, err = r.env.CSVInside("CITY", constants.SurfsharkCityChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable CITY: %w", err)
	}

	settings.ServerSelection.Hostnames, err = r.env.CSVInside("SERVER_HOSTNAME",
		constants.SurfsharkHostnameChoices(servers))
	if err != nil {
		return fmt.Errorf("environment variable SERVER_HOSTNAME: %w", err)
	}

	regionChoices := constants.SurfsharkRegionChoices(servers)
	regionChoices = append(regionChoices, constants.SurfsharkRetroLocChoices(servers)...)
	settings.ServerSelection.Regions, err = r.env.CSVInside("REGION", regionChoices)
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	// Retro compatibility
	// TODO remove in v4
	settings.ServerSelection = surfsharkRetroRegion(settings.ServerSelection)

	settings.ServerSelection.MultiHopOnly, err = r.env.YesNo("MULTIHOP_ONLY", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable MULTIHOP_ONLY: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolOnly(r)
}

func surfsharkRetroRegion(selection ServerSelection) (
	updatedSelection ServerSelection) {
	locationData := constants.SurfsharkLocationData()

	retroToLocation := make(map[string]models.SurfsharkLocationData, len(locationData))
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
