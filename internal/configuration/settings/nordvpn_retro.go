package settings

// Retro-compatibility because SERVER_REGIONS changed to SERVER_COUNTRIES
// and SERVER_REGIONS is now the continent field for servers.
// TODO v4 remove.
func nordvpnRetroRegion(selection ServerSelection, validRegions, validCountries []string) (
	updatedSelection ServerSelection,
) {
	validRegionsMap := stringSliceToMap(validRegions)
	validCountriesMap := stringSliceToMap(validCountries)

	updatedSelection = selection.copy()
	updatedSelection.Regions = make([]string, 0, len(selection.Regions))
	for _, region := range selection.Regions {
		_, isValid := validRegionsMap[region]
		if isValid {
			updatedSelection.Regions = append(updatedSelection.Regions, region)
			continue
		}

		_, isValid = validCountriesMap[region]
		if !isValid {
			// Region is not valid for the country or region
			// just leave it to the validation to fail it later
			continue
		}

		// Region is not valid for a region, but is a valid country
		// Handle retro-compatibility and transfer the value to the
		// country field.
		updatedSelection.Countries = append(updatedSelection.Countries, region)
	}

	return updatedSelection
}

func stringSliceToMap(slice []string) (m map[string]struct{}) {
	m = make(map[string]struct{}, len(slice))
	for _, s := range slice {
		m[s] = struct{}{}
	}
	return m
}
