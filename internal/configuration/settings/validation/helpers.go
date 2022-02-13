package validation

import "sort"

func makeUnique(choices []string) (uniqueChoices []string) {
	seen := make(map[string]struct{}, len(choices))
	uniqueChoices = make([]string, 0, len(uniqueChoices))

	for _, choice := range choices {
		if _, ok := seen[choice]; ok {
			continue
		}
		seen[choice] = struct{}{}

		uniqueChoices = append(uniqueChoices, choice)
	}

	sort.Slice(uniqueChoices, func(i, j int) bool {
		return uniqueChoices[i] < uniqueChoices[j]
	})

	return uniqueChoices
}
