package mux

import (
	"fmt"

	"github.com/qdm12/gluetun/internal/configuration/settings"
	"github.com/qdm12/gluetun/internal/configuration/sources"
)

var _ sources.Source = (*Reader)(nil)

type Reader struct {
	sources []sources.Source
}

func New(sources ...sources.Source) *Reader {
	return &Reader{
		sources: sources,
	}
}

// Read reads the settings for each source, merging unset fields
// with field set by the next source.
// It then set defaults to remaining unset fields.
func (r *Reader) Read() (settings settings.Settings, err error) {
	for _, source := range r.sources {
		settingsFromSource, err := source.Read()
		if err != nil {
			return settings, fmt.Errorf("cannot read from source %T: %w", source, err)
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
func (r *Reader) ReadHealth() (settings settings.Health, err error) {
	for _, source := range r.sources {
		settingsFromSource, err := source.ReadHealth()
		if err != nil {
			return settings, fmt.Errorf("cannot read from source %T: %w", source, err)
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
