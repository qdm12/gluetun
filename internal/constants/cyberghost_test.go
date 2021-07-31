package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CyberghostGroupChoices(t *testing.T) {
	t.Parallel()

	expected := []string{"Premium TCP Asia", "Premium TCP Europe",
		"Premium TCP USA", "Premium UDP Asia", "Premium UDP Europe",
		"Premium UDP USA"}
	choices := CyberghostGroupChoices()

	assert.Equal(t, expected, choices)
}
