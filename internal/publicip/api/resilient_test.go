package api

import (
	"testing"

	"github.com/qdm12/gluetun/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_GetMostPopularResult(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		input    []models.PublicIP
		expected models.PublicIP
	}{
		"exact_matches": {
			input: []models.PublicIP{
				{Country: "France", City: "Paris"},
				{Country: "USA", City: "New York"},
				{Country: "France", City: "Paris"},
			},
			expected: models.PublicIP{Country: "France", City: "Paris"},
		},
		"fuzzy_country_matching": {
			input: []models.PublicIP{
				{Country: "Germany", Region: "Bavaria", City: "Munich"},
				{Country: "Germani", Region: "Bavaria", City: "Munich"},
				{Country: "France", Region: "IDF", City: "Paris"},
			},
			expected: models.PublicIP{Country: "Germany", Region: "Bavaria", City: "Munich"},
		},
		"hierarchy_priority": {
			input: []models.PublicIP{
				{Country: "Italy", Region: "Sicily", City: "Syracuse"},
				{Country: "Italy", Region: "Sicily", City: "Syracuse"},
				{Country: "USA", Region: "New York", City: "Syracuse"},
				{Country: "Italy", Region: "Sicily", City: "Syracuse"},
			},
			expected: models.PublicIP{Country: "Italy", Region: "Sicily", City: "Syracuse"},
		},
		"normalization_check": {
			input: []models.PublicIP{
				{Country: "Canada", City: "Montréal"},
				{Country: "Canada", City: "Montreal "},
				{Country: "UK", City: "London"},
			},
			expected: models.PublicIP{Country: "Canada", City: "Montréal"},
		},
		"all_different": {
			input: []models.PublicIP{
				{Country: "Canada", City: "Montréal"},
				{Country: "US", City: "New York"},
				{Country: "UK", City: "London"},
			},
			expected: models.PublicIP{Country: "US", City: "New York"},
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := getMostPopularResult(testCase.input)

			assert.Equal(t, testCase.expected.Country, result.Country)
			assert.Equal(t, testCase.expected.Region, result.Region)
			assert.Equal(t, testCase.expected.City, result.City)
		})
	}
}
