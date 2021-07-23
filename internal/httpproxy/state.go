package httpproxy

import (
	"sync"

	"github.com/qdm12/gluetun/internal/configuration"
	"github.com/qdm12/gluetun/internal/models"
)

type state struct {
	status     models.LoopStatus
	settings   configuration.HTTPProxy
	statusMu   sync.RWMutex
	settingsMu sync.RWMutex
}
