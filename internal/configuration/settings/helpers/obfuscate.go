package helpers

func ObfuscateWireguardKey(fullKey string) (obfuscatedKey string) {
	const minKeyLength = 10
	if len(fullKey) < minKeyLength {
		return "(too short)"
	}

	lastIndex := len(fullKey) - 1
	return fullKey[0:2] + "..." + fullKey[lastIndex-2:]
}

func ObfuscatePassword(password string) (obfuscatedPassword string) {
	if password != "" {
		return "[set]"
	}
	return "[not set]"
}

func ObfuscateData(data string) (obfuscated string) {
	if data != "" {
		return "[set]"
	}
	return "[not set]"
}
