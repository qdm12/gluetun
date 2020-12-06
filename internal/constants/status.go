package constants

import (
	"github.com/qdm12/gluetun/internal/models"
)

const (
	Starting  models.LoopStatus = "starting"
	Running   models.LoopStatus = "running"
	Stopping  models.LoopStatus = "stopping"
	Stopped   models.LoopStatus = "stopped"
	Crashed   models.LoopStatus = "crashed"
	Completed models.LoopStatus = "completed"
)
