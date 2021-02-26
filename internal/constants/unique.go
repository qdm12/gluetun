package constants

import "sort"

func makeUnique(slice []string) (uniques []string) {
	set := make(map[string]struct{}, len(slice))
	for _, element := range slice {
		set[element] = struct{}{}
	}

	uniques = make([]string, 0, len(set))
	for element := range set {
		uniques = append(uniques, element)
	}

	sort.Slice(uniques, func(i, j int) bool {
		return uniques[i] < uniques[j]
	})

	return uniques
}
