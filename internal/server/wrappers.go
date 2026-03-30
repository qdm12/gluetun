package server

import (
	"errors"
	"fmt"

	"github.com/qdm12/gluetun/internal/constants"
	"github.com/qdm12/gluetun/internal/models"
)

type statusWrapper struct {
	Status string `json:"status"`
}

var errInvalidStatus = errors.New("invalid status")

func (sw *statusWrapper) getStatus() (status models.LoopStatus, err error) {
	status = models.LoopStatus(sw.Status)
	switch status {
	case constants.Stopped, constants.Running:
		return status, nil
	default:
		return "", fmt.Errorf("%w: %s: possible values are: %s, %s",
			errInvalidStatus, sw.Status, constants.Stopped, constants.Running)
	}
}

type portsWrapper struct {
	Port  uint16   `json:"port"` // TODO v4 remove
	Ports []uint16 `json:"ports"`
}

type outcomeWrapper struct {
	Outcome string `json:"outcome"`
}
