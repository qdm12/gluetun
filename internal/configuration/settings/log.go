package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
	"github.com/qdm12/log"
)

// Log contains settings to configure the logger.
type Log struct {
	// Level is the log level of the logger.
	// It cannot be nil in the internal state.
	Level *log.Level
}

func (l Log) validate() (err error) {
	return nil
}

func (l *Log) copy() (copied Log) {
	return Log{
		Level: helpers.CopyLogLevelPtr(l.Level),
	}
}

// mergeWith merges the other settings into any
// unset field of the receiver settings object.
func (l *Log) mergeWith(other Log) {
	l.Level = helpers.MergeWithLogLevel(l.Level, other.Level)
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (l *Log) overrideWith(other Log) {
	l.Level = helpers.OverrideWithLogLevel(l.Level, other.Level)
}

func (l *Log) setDefaults() {
	l.Level = helpers.DefaultLogLevel(l.Level, log.LevelInfo)
}

func (l Log) String() string {
	return l.toLinesNode().String()
}

func (l Log) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Log settings:")
	node.Appendf("Log level: %s", l.Level.String())
	return node
}
