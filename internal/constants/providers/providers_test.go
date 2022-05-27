package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_All(t *testing.T) {
	t.Parallel()

	all := All()
	assert.NotContains(t, all, Custom)
	assert.NotEmpty(t, all)
}

func Test_AllWithCustom(t *testing.T) {
	t.Parallel()

	all := AllWithCustom()
	assert.Contains(t, all, Custom)
	assert.Len(t, all, len(All())+1)
}
