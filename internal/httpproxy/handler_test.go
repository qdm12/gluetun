package httpproxy

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_returnRedirect(t *testing.T) {
	t.Parallel()

	err := returnRedirect(nil, nil)

	assert.Equal(t, http.ErrUseLastResponse, err)
}
