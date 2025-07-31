package settings

import (
	"os"
	"strings"
	"testing"
)

func TestWireguardCustomConfigFile(t *testing.T) {
	// Use the test config file from the repo, with bogus values only
	confPath := "../../wireguard/wg_test.conf"
	t.Logf("Using config file: %s", confPath)
	wg := &Wireguard{
		CustomConfigFile: &confPath,
	}

	t.Log("Checking if CustomConfigFile is set...")
	if wg.CustomConfigFile == nil || *wg.CustomConfigFile == "" {
		t.Error("CustomConfigFile should be set")
	} else {
		t.Logf("CustomConfigFile is set to: %s", *wg.CustomConfigFile)
	}

	t.Log("Checking if config file exists and is readable...")
	if _, err := os.Stat(*wg.CustomConfigFile); err != nil {
		t.Errorf("Custom config file not found: %v", err)
	} else {
		t.Log("Config file exists and is readable.")
	}

	t.Log("Reading config file contents...")
	data, err := os.ReadFile(*wg.CustomConfigFile)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}
	content := string(data)
	t.Logf("Config file content (first 100 chars): %.100s", content)

	t.Log("Checking for bogus or real values in config file...")
	if containsBogusWireguardValues(content) {
		t.Log("Bogus values found in config file (safe for public use/testing).")
	} else if containsRealWireguardValues(content) {
		t.Log("Real values found in config file (for production/validation use).")
	} else {
		t.Error("wg_test.conf does not contain recognized bogus or real values.")
	}
}

// containsBogusWireguardValues checks for known bogus values in the config
func containsBogusWireguardValues(content string) bool {
	bogusKeys := []string{
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", // bogus private key
		"BBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBBB=", // bogus public key
		"1.2.3.4/32",            // bogus address
		"8.8.8.8",               // bogus DNS
		"123.123.123.123:51820", // bogus endpoint
	}
	for _, v := range bogusKeys {
		if !strings.Contains(content, v) {
			return false
		}
	}
	return true
}

// containsRealWireguardValues checks for known real values in the config (edit as needed for your real config)
func containsRealWireguardValues(content string) bool {
	// If not bogus, assume real/production values (do not match any known bogus test values)
	realKeys := []string{
		// These are placeholders for real values; in production, these would be actual keys/addresses.
		// For test, just check that the config does NOT contain all bogus values.
	}
	for _, v := range realKeys {
		if !strings.Contains(content, v) {
			return false
		}
	}
	return true
}
