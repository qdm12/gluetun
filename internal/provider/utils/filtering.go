package utils

import "strings"

func FilterByPossibilities(value string, possibilities []string) (filtered bool) {
	if len(possibilities) == 0 {
		return false
	}
	for _, possibility := range possibilities {
		if strings.EqualFold(value, possibility) {
			return false
		}
	}
	return true
}
