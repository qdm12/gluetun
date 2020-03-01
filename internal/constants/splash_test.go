package constants

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AnnoucementExpiration(t *testing.T) {
	t.Parallel()
	if len(AnnoucementExpiration) == 0 {
		return
	}
	_, err := time.Parse("2006-01-02", AnnoucementExpiration)
	assert.NoError(t, err)
}
