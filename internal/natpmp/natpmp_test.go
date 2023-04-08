package natpmp

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()

	expectedClient := &Client{
		serverPort:   5351,
		initialRetry: 250 * time.Millisecond,
		maxRetries:   9,
	}
	client := New()
	assert.Equal(t, expectedClient, client)
}
