package wireguard

import (
	"os"
	"testing"
)

func TestParseConfFile(t *testing.T) {
	confPath := "./wg_test.conf"
	if _, err := os.Stat(confPath); err != nil {
		t.Fatalf("test config file not found: %v", err)
	}

	parsed, err := ParseConfFile(confPath)
	if err != nil {
		t.Fatalf("ParseConfFile failed: %v", err)
	}

	provider := DetectProvider(parsed)
	t.Logf("Detected provider: %s", provider)

	switch provider {
	case "protonvpn":
		if parsed.Interface["PrivateKey"] == "" {
			t.Error("ProtonVPN: PrivateKey not parsed correctly")
		}
		if parsed.Peer["PublicKey"] == "" {
			t.Error("ProtonVPN: PublicKey not parsed correctly")
		}
		if parsed.Peer["Endpoint"] == "" {
			t.Error("ProtonVPN: Endpoint not parsed correctly")
		}
	case "mullvad":
		// Add Mullvad-specific checks here
		if parsed.Interface["PrivateKey"] == "" {
			t.Error("Mullvad: PrivateKey not parsed correctly")
		}
		// ...
	default:
		t.Log("Unknown or generic provider, basic checks only.")
		if parsed.Interface["PrivateKey"] == "" {
			t.Error("PrivateKey not parsed correctly")
		}
	}
}
