package internal

import (
	"context"
	"fmt"
)

func ProtonVPNTest(ctx context.Context, logger Logger) error {
	expectedSecrets := []string{
		"Wireguard private key",
	}
	secrets, err := readSecrets(ctx, expectedSecrets, logger)
	if err != nil {
		return fmt.Errorf("reading secrets: %w", err)
	}

	env := []string{
		"VPN_SERVICE_PROVIDER=protonvpn",
		"VPN_TYPE=wireguard",
		"LOG_LEVEL=debug",
		"SERVER_COUNTRIES=United States",
		"WIREGUARD_PRIVATE_KEY=" + secrets[0],
	}
	return simpleTest(ctx, env, logger)
}
