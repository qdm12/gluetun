package auth

func andStrings(strings []string) (result string) {
	return joinStrings(strings, "and")
}

func orStrings(strings []string) (result string) {
	return joinStrings(strings, "or")
}

func joinStrings(strings []string, lastJoin string) (result string) {
	if len(strings) == 0 {
		return ""
	}

	result = strings[0]
	for i := 1; i < len(strings); i++ {
		if i < len(strings)-1 {
			result += ", " + strings[i]
		} else {
			result += " " + lastJoin + " " + strings[i]
		}
	}

	return result
}
