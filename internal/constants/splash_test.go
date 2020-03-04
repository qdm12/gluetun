package constants

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AnnouncementExpiration(t *testing.T) {
	t.Parallel()
	if len(AnnouncementExpiration) == 0 {
		return
	}
	_, err := time.Parse("2006-01-02", AnnouncementExpiration)
	assert.NoError(t, err)
}
