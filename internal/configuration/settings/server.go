package settings

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/qdm12/gluetun/internal/server/middlewares/auth"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
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
	// AuthFilePath is the path to the file containing the authentication
	// configuration for the middleware.
	// It cannot be empty in the internal state and defaults to
	// /gluetun/auth/config.toml.
	AuthFilePath string
	// Auth contains settings for the authentication middleware.
	// These are parsed from a configuration file specified by
	// AuthFilePath.
	Auth auth.Settings
}

func (c ControlServer) validate() (err error) {
	_, portStr, err := net.SplitHostPort(*c.Address)
	if err != nil {
		return fmt.Errorf("listening address is not valid: %w", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("listening port it not valid: %w", err)
	}

	uid := os.Getuid()
	const maxPrivilegedPort = 1023
	if uid != 0 && port != 0 && port <= maxPrivilegedPort {
		return fmt.Errorf("%w: %d when running with user ID %d",
			ErrControlServerPrivilegedPort, port, uid)
	}

	err = c.Auth.Validate()
	if err != nil {
		return fmt.Errorf("validating authentication middleware: %w", err)
	}

	return nil
}

func (c *ControlServer) copy() (copied ControlServer) {
	return ControlServer{
		Address: gosettings.CopyPointer(c.Address),
		Log:     gosettings.CopyPointer(c.Log),
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (c *ControlServer) overrideWith(other ControlServer) {
	c.Address = gosettings.OverrideWithPointer(c.Address, other.Address)
	c.Log = gosettings.OverrideWithPointer(c.Log, other.Log)
	c.AuthFilePath = gosettings.OverrideWithComparable(c.AuthFilePath, other.AuthFilePath)
	c.Auth.OverrideWith(other.Auth)
}

func (c *ControlServer) setDefaults() {
	c.Address = gosettings.DefaultPointer(c.Address, ":8000")
	c.Log = gosettings.DefaultPointer(c.Log, true)
	c.AuthFilePath = gosettings.DefaultComparable(c.AuthFilePath, "/gluetun/auth/config.toml")
	c.Auth.SetDefaults()
}

func (c ControlServer) String() string {
	return c.toLinesNode().String()
}

func (c ControlServer) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Control server settings:")
	node.Appendf("Listening address: %s", *c.Address)
	node.Appendf("Logging: %s", gosettings.BoolToYesNo(c.Log))
	node.Appendf("Authentication file path: %s", c.AuthFilePath)
	node.AppendNode(c.Auth.ToLinesNode())
	return node
}

func (c *ControlServer) read(r *reader.Reader) (err error) {
	c.Log, err = r.BoolPtr("HTTP_CONTROL_SERVER_LOG")
	if err != nil {
		return err
	}

	c.Address = r.Get("HTTP_CONTROL_SERVER_ADDRESS")

	c.AuthFilePath = r.String("HTTP_CONTROL_SERVER_AUTH_CONFIG_FILEPATH")
	if c.AuthFilePath != "" {
		c.Auth, err = auth.Read(c.AuthFilePath)
		if err != nil {
			return fmt.Errorf("reading authentication middleware settings: %w", err)
		}
	}

	return nil
}
