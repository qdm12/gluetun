package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
)

// System contains settings to configure system related elements.
type System struct {
	PUID     *uint16 `json:"puid,omitempty"`
	PGID     *uint16 `json:"pgid,omitempty"`
	Timezone string  `json:"timezone,omitempty"`
}

// Validate validates System settings.
func (s System) validate() (err error) {
	return nil
}

func (s *System) copy() (copied System) {
	return System{
		PUID:     helpers.CopyUint16Ptr(s.PUID),
		PGID:     helpers.CopyUint16Ptr(s.PGID),
		Timezone: s.Timezone,
	}
}

func (s *System) mergeWith(other System) {
	s.PUID = helpers.MergeWithUint16(s.PUID, other.PUID)
	s.PGID = helpers.MergeWithUint16(s.PGID, other.PGID)
	s.Timezone = helpers.MergeWithString(s.Timezone, other.Timezone)
}

func (s *System) overrideWith(other System) {
	s.PUID = helpers.OverrideWithUint16(s.PUID, other.PUID)
	s.PGID = helpers.OverrideWithUint16(s.PGID, other.PGID)
	s.Timezone = helpers.OverrideWithString(s.Timezone, other.Timezone)
}

func (s *System) setDefaults() {
	const defaultID = 1000
	s.PUID = helpers.DefaultUint16(s.PUID, defaultID)
	s.PGID = helpers.DefaultUint16(s.PGID, defaultID)
}
