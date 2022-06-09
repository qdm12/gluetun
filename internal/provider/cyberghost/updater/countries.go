package updater

func mergeCountryCodes(base, extend map[string]string) (merged map[string]string) {
	merged = make(map[string]string, len(base))
	for countryCode, region := range base {
		merged[countryCode] = region
	}
	for countryCode := range base {
		delete(extend, countryCode)
	}
	for countryCode, region := range extend {
		merged[countryCode] = region
	}
	return merged
}
