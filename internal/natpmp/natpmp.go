package natpmp

import (
	"time"
)

// Client is a NAT-PMP protocol client.
type Client struct {
	serverPort   uint16
	initialRetry time.Duration
	maxRetries   uint
}

// New creates a new NAT-PMP client.
func New() (client *Client) {
	const natpmpPort = 5351

	// Parameters described in https://www.ietf.org/rfc/rfc6886.html#section-3.1
	const initialRetry = 250 * time.Millisecond
	const maxTries = 9 // 64 seconds
	return &Client{
		serverPort:   natpmpPort,
		initialRetry: initialRetry,
		maxRetries:   maxTries,
	}
}
