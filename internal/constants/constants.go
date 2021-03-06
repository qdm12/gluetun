// Package constants defines constants shared throughout the program.
// It also defines constant maps and slices using functions.
package constants

import "sort"

func makeChoicesUnique(choices []string) []string {
	uniqueChoices := map[string]struct{}{}
	for _, choice := range choices {
		uniqueChoices[choice] = struct{}{}
	}

	uniqueChoicesSlice := make([]string, len(uniqueChoices))
	i := 0
	for choice := range uniqueChoices {
		uniqueChoicesSlice[i] = choice
		i++
	}

	sort.Slice(uniqueChoicesSlice, func(i, j int) bool {
		return uniqueChoicesSlice[i] < uniqueChoicesSlice[j]
	})

	return uniqueChoicesSlice
}
