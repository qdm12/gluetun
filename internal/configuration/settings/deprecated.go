package settings

import (
	"slices"

	"github.com/qdm12/gosettings/reader"
	"golang.org/x/exp/maps"
)

func readObsolete(r *reader.Reader) (warnings []string) {
	keyToMessage := map[string]string{
		"DOT_VERBOSITY":           "DOT_VERBOSITY is obsolete, use LOG_LEVEL instead.",
		"DOT_VERBOSITY_DETAILS":   "DOT_VERBOSITY_DETAILS is obsolete because it was specific to Unbound.",
		"DOT_VALIDATION_LOGLEVEL": "DOT_VALIDATION_LOGLEVEL is obsolete because DNSSEC validation is not implemented.",
	}
	sortedKeys := maps.Keys(keyToMessage)
	slices.Sort(sortedKeys)
	warnings = make([]string, 0, len(keyToMessage))
	for _, key := range sortedKeys {
		if r.Get(key) != nil {
			warnings = append(warnings, keyToMessage[key])
		}
	}
	return warnings
}
