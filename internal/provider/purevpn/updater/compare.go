package updater

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// comparePlaceNames returns true if strings are within 1 edit
// distance after normalization
func comparePlaceNames(a, b string) bool {
	normA := normalize(a)
	normB := normalize(b)
	return normA == normB || levenshteinDistance(normA, normB) <= 1
}

// normalize removes accents, trims space, and lowercases the string
func normalize(s string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(transformer, s)
	if err != nil {
		panic(err)
	}
	return strings.ToLower(strings.TrimSpace(result))
}

// levenshteinDistance calculates the edit distance
// between two strings a and b.
func levenshteinDistance(a, b string) int {
	switch {
	case len(a) == 0:
		return len(b)
	case len(b) == 0:
		return len(a)
	}

	column := make([]int, len(b)+1)
	for i := 0; i <= len(b); i++ {
		column[i] = i
	}

	for i := 1; i <= len(a); i++ {
		column[0] = i
		lastValue := i - 1
		for j := 1; j <= len(b); j++ {
			oldValue := column[j]
			cost := 0
			if a[i-1] != b[j-1] {
				cost = 1
			}
			column[j] = min(column[j]+1, min(column[j-1]+1, lastValue+cost))
			lastValue = oldValue
		}
	}
	return column[len(b)]
}
