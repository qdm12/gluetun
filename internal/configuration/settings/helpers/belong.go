package helpers

func IsOneOf[T comparable](value T, choices ...T) (ok bool) {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}
