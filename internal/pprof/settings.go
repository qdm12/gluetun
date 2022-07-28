package pprof

import (
	"errors"
	"time"

	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gluetun/internal/httpserver"
	"github.com/qdm12/gotree"
)

// Settings are the settings for the Pprof service.
type Settings struct {
	// Enabled can be false or true.
	// It defaults to false.
	Enabled *bool
	// See runtime.SetBlockProfileRate
	// Set to 0 to disable profiling.
	BlockProfileRate int
	// See runtime.SetMutexProfileFraction
	// Set to 0 to disable profiling.
	MutexProfileRate int
	// HTTPServer contains settings to configure
	// the HTTP server serving pprof data.
	HTTPServer httpserver.Settings
}

func (s *Settings) SetDefaults() {
	s.Enabled = helpers.DefaultBool(s.Enabled, false)
	s.HTTPServer.Address = helpers.DefaultString(s.HTTPServer.Address, "localhost:6060")
	const defaultReadTimeout = 5 * time.Minute // for CPU profiling
	s.HTTPServer.ReadTimeout = helpers.DefaultDuration(s.HTTPServer.ReadTimeout, defaultReadTimeout)
	s.HTTPServer.SetDefaults()
}

func (s Settings) Copy() (copied Settings) {
	return Settings{
		Enabled:          helpers.CopyBoolPtr(s.Enabled),
		BlockProfileRate: s.BlockProfileRate,
		MutexProfileRate: s.MutexProfileRate,
		HTTPServer:       s.HTTPServer.Copy(),
	}
}

func (s *Settings) MergeWith(other Settings) {
	s.Enabled = helpers.MergeWithBool(s.Enabled, other.Enabled)
	s.BlockProfileRate = helpers.MergeWithInt(s.BlockProfileRate, other.BlockProfileRate)
	s.MutexProfileRate = helpers.MergeWithInt(s.MutexProfileRate, other.MutexProfileRate)
	s.HTTPServer.MergeWith(other.HTTPServer)
}

func (s *Settings) OverrideWith(other Settings) {
	s.Enabled = helpers.OverrideWithBool(s.Enabled, other.Enabled)
	s.BlockProfileRate = helpers.OverrideWithInt(s.BlockProfileRate, other.BlockProfileRate)
	s.MutexProfileRate = helpers.OverrideWithInt(s.MutexProfileRate, other.MutexProfileRate)
	s.HTTPServer.OverrideWith(other.HTTPServer)
}

var (
	ErrBlockProfileRateNegative = errors.New("block profile rate cannot be negative")
	ErrMutexProfileRateNegative = errors.New("mutex profile rate cannot be negative")
)

func (s Settings) Validate() (err error) {
	if s.BlockProfileRate < 0 {
		return ErrBlockProfileRateNegative
	}

	if s.MutexProfileRate < 0 {
		return ErrMutexProfileRateNegative
	}

	return s.HTTPServer.Validate()
}

func (s Settings) ToLinesNode() (node *gotree.Node) {
	if !*s.Enabled {
		return nil
	}

	node = gotree.New("Pprof settings:")

	if s.BlockProfileRate > 0 {
		node.Appendf("Block profile rate: %d", s.BlockProfileRate)
	}

	if s.MutexProfileRate > 0 {
		node.Appendf("Mutex profile rate: %d", s.MutexProfileRate)
	}

	node.AppendNode(s.HTTPServer.ToLinesNode())

	return node
}

func (s Settings) String() string {
	return s.ToLinesNode().String()
}
