package merge

import (
	"fmt"
	"strings"

	"github.com/qdm12/gluetun/internal/configuration/settings"
)

type ConfigSource interface {
	Read() (settings settings.Settings, err error)
	ReadHealth() (settings settings.Health, err error)
	String() string
}

type Source struct {
	sources []ConfigSource
}

func New(sources ...ConfigSource) *Source {
	return &Source{
		sources: sources,
	}
}

func (s *Source) String() string {
	sources := make([]string, len(s.sources))
	for i := range s.sources {
		sources[i] = s.sources[i].String()
	}
	return strings.Join(sources, ", ")
}

// Read reads the settings for each source, merging unset fields
// with field set by the next source.
// It then set defaults to remaining unset fields.
func (s *Source) Read() (settings settings.Settings, err error) {
	for _, source := range s.sources {
		settingsFromSource, err := source.Read()
		if err != nil {
			return settings, fmt.Errorf("reading from %s: %w", source, err)
		}
		settings.MergeWith(settingsFromSource)
	}
	settings.SetDefaults()
	return settings, nil
}

// ReadHealth reads the health settings for each source, merging unset fields
// with field set by the next source.
// It then set defaults to remaining unset fields, and validate
// all the fields.
func (s *Source) ReadHealth() (settings settings.Health, err error) {
	for _, source := range s.sources {
		settingsFromSource, err := source.ReadHealth()
		if err != nil {
			return settings, fmt.Errorf("reading from %s: %w", source, err)
		}
		settings.MergeWith(settingsFromSource)
	}
	settings.SetDefaults()

	err = settings.Validate()
	if err != nil {
		return settings, err
	}

	return settings, nil
}
