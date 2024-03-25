package settings

import (
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
	"github.com/qdm12/log"
)

// Log contains settings to configure the logger.
type Log struct {
	// Level is the log level of the logger.
	// It cannot be empty in the internal state.
	Level string
}

func (l Log) validate() (err error) {
	_, err = log.ParseLevel(l.Level)
	if err != nil {
		return fmt.Errorf("level: %w", err)
	}
	return nil
}

func (l *Log) copy() (copied Log) {
	return Log{
		Level: l.Level,
	}
}

// overrideWith overrides fields of the receiver
// settings object with any field set in the other
// settings.
func (l *Log) overrideWith(other Log) {
	l.Level = gosettings.OverrideWithComparable(l.Level, other.Level)
}

func (l *Log) setDefaults() {
	l.Level = gosettings.DefaultComparable(l.Level, log.LevelInfo.String())
}

func (l Log) String() string {
	return l.toLinesNode().String()
}

func (l Log) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Log settings:")
	node.Appendf("Log level: %s", l.Level)
	return node
}

func (l *Log) read(r *reader.Reader) (err error) {
	l.Level = r.String("LOG_LEVEL")
	return nil
}
