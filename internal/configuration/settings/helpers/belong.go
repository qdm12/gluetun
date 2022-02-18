package helpers

import (
	"errors"
	"fmt"
	"strings"
)

func IsOneOf(value string, choices ...string) (ok bool) {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}

var (
	ErrNoChoice      = errors.New("one or more values is set but there is no possible value available")
	ErrValueNotOneOf = errors.New("value is not one of the possible choices")
)

func AreAllOneOf(values, choices []string) (err error) {
	if len(values) > 0 && len(choices) == 0 {
		return ErrNoChoice
	}

	set := make(map[string]struct{}, len(choices))
	for _, choice := range choices {
		choice = strings.ToLower(choice)
		set[choice] = struct{}{}
	}

	for _, value := range values {
		_, ok := set[value]
		if !ok {
			return fmt.Errorf("%w: value %q, choices available are %s",
				ErrValueNotOneOf, value, strings.Join(choices, ", "))
		}
	}

	return nil
}

func Uint16IsOneOf(port uint16, choices []uint16) (ok bool) {
	for _, choice := range choices {
		if port == choice {
			return true
		}
	}
	return false
}
