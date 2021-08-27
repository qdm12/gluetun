package configuration

import (
	"fmt"

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
	regions, err := r.env.CSVInside("REGION", regionChoices)
	if err != nil {
		return fmt.Errorf("environment variable REGION: %w", err)
	}

	// Retro compatibility
	// TODO remove in v4
	for i, region := range regions {
		locationData, isRetro :=
			surfsharkConvertRetroLoc(region)
		if !isRetro {
			continue
		}

		regions[i] = locationData.Region
		settings.ServerSelection.Countries = append(settings.ServerSelection.Countries, locationData.Country)
		if locationData.City != "" { // city is empty for some servers
			settings.ServerSelection.Cities = append(settings.ServerSelection.Cities, locationData.City)
		}
		settings.ServerSelection.Hostnames = append(settings.ServerSelection.Hostnames, locationData.Hostname)
	}

	settings.ServerSelection.MultiHopOnly, err = r.env.YesNo("MULTIHOP_ONLY", params.Default("no"))
	if err != nil {
		return fmt.Errorf("environment variable MULTIHOP_ONLY: %w", err)
	}

	return settings.ServerSelection.OpenVPN.readProtocolOnly(r.env)
}

// TODO remove in v4.
func surfsharkConvertRetroLoc(retroLoc string) (
	locationData models.SurfsharkLocationData, isRetro bool) {
	for _, data := range constants.SurfsharkLocationData() {
		if retroLoc == data.RetroLoc {
			return data, true
		}
	}
	return locationData, false
}
