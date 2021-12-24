package settings

import (
	"fmt"
	"os"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
)

// ControlServer contains settings to customize the control server operation.
type ControlServer struct {
	// Port is the listening port to use.
	// It can be set to 0 to bind to a random port.
	// It cannot be nil in the internal state.
	// TODO change to address
	Port *uint16 `json:"port,omitempty"`
	// Log can be true or false to enable logging on requests.
	// It cannot be nil in the internal state.
	Log *bool `json:"log,omitempty"`
}

func (c ControlServer) validate() (err error) {
	uid := os.Getuid()
	const maxPrivilegedPort uint16 = 1023
	if uid != 0 && *c.Port <= maxPrivilegedPort {
		return fmt.Errorf("%w: %d when running with user ID %d",
			ErrControlServerPrivilegedPort, *c.Port, uid)
	}

	return nil
}

func (c *ControlServer) copy() (copied ControlServer) {
	return ControlServer{
		Port: helpers.CopyUint16Ptr(c.Port),
		Log:  helpers.CopyBoolPtr(c.Log),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (c *ControlServer) mergeWith(other ControlServer) {
	c.Port = helpers.MergeWithUint16(c.Port, other.Port)
	c.Log = helpers.MergeWithBool(c.Log, other.Log)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (c *ControlServer) overrideWith(other ControlServer) {
	c.Port = helpers.MergeWithUint16(c.Port, other.Port)
	c.Log = helpers.MergeWithBool(c.Log, other.Log)
}

func (c *ControlServer) setDefaults() {
	const defaultPort = 8000
	c.Port = helpers.DefaultUint16(c.Port, defaultPort)
	c.Log = helpers.DefaultBool(c.Log, true)
}
