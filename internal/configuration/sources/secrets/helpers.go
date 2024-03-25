package secrets

func strPtrToStringIsSet(ptr *string) (s string, isSet bool) {
	if ptr == nil {
		return "", false
	}
	return *ptr, true
}
