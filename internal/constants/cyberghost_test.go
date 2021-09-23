package constants

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/qdm12/gluetun/internal/storage"
	"github.com/qdm12/golibs/logging/mock_logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CyberghostGroupChoices(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)

	logger := mock_logging.NewMockLogger(ctrl)
	logger.EXPECT().Info(gomock.Any())

	storage, err := storage.New(logger, "")
	require.NoError(t, err)

	servers := storage.GetServers()

	expected := []string{"NoSpy TCP Europe", "NoSpy UDP Europe",
		"Premium TCP Europe", "Premium UDP Europe"}
	choices := CyberghostGroupChoices(servers.GetCyberghost())

	assert.Equal(t, expected, choices)
}
