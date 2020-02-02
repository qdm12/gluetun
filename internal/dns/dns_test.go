package dns

import (
	"testing"

	"github.com/qdm12/golibs/files"
	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/network"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewConfigurator(t *testing.T) {
	t.Parallel()
	logger, err := logging.NewEmptyLogging()
	require.NoError(t, err)
	client := network.NewClient(0)
	fileManager := files.NewFileManager()
	var c Configurator
	c = NewConfigurator(logger, client, fileManager)
	assert.Equal(t, &configurator{
		logger:      logger,
		client:      client,
		fileManager: fileManager,
	}, c)
}
