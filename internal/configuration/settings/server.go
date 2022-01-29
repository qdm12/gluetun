package settings

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
)

// ControlServer contains settings to customize the control server operation.
type ControlServer struct {
	// Address is the listening address to use.
	// It cannot be nil in the internal state.
	Address *string
	// Log can be true or false to enable logging on requests.
	// It cannot be nil in the internal state.
	Log *bool
}

func (c ControlServer) validate() (err error) {
	_, portStr, err := net.SplitHostPort(*c.Address)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrControlServerAddress, err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrControlServerPort, err)
	}

	uid := os.Getuid()
	const maxPrivilegedPort = 1023
	if uid != 0 && port != 0 && port <= maxPrivilegedPort {
		return fmt.Errorf("%w: %d when running with user ID %d",
			ErrControlServerPrivilegedPort, port, uid)
	}

	return nil
}

func (c *ControlServer) copy() (copied ControlServer) {
	return ControlServer{
		Address: helpers.CopyStringPtr(c.Address),
		Log:     helpers.CopyBoolPtr(c.Log),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (c *ControlServer) mergeWith(other ControlServer) {
	c.Address = helpers.MergeWithStringPtr(c.Address, other.Address)
	c.Log = helpers.MergeWithBool(c.Log, other.Log)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (c *ControlServer) overrideWith(other ControlServer) {
	c.Address = helpers.OverrideWithStringPtr(c.Address, other.Address)
	c.Log = helpers.OverrideWithBool(c.Log, other.Log)
}

func (c *ControlServer) setDefaults() {
	c.Address = helpers.DefaultStringPtr(c.Address, ":8000")
	c.Log = helpers.DefaultBool(c.Log, true)
}

func (c ControlServer) String() string {
	return c.toLinesNode().String()
}

func (c ControlServer) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Control server settings:")
	node.Appendf("Listening address: %s", *c.Address)
	node.Appendf("Logging: %s", helpers.BoolPtrToYesNo(c.Log))
	return node
}
