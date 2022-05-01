package settings

import (
	"github.com/qdm12/gluetun/internal/configuration/settings/helpers"
	"github.com/qdm12/gotree"
)

// System contains settings to configure system related elements.
type System struct {
	PUID     *uint32
	PGID     *uint32
	Timezone string
}

// Validate validates System settings.
func (s System) validate() (err error) {
	return nil
}

func (s *System) copy() (copied System) {
	return System{
		PUID:     helpers.CopyUint32Ptr(s.PUID),
		PGID:     helpers.CopyUint32Ptr(s.PGID),
		Timezone: s.Timezone,
	}
}

func (s *System) mergeWith(other System) {
	s.PUID = helpers.MergeWithUint32(s.PUID, other.PUID)
	s.PGID = helpers.MergeWithUint32(s.PGID, other.PGID)
	s.Timezone = helpers.MergeWithString(s.Timezone, other.Timezone)
}

func (s *System) overrideWith(other System) {
	s.PUID = helpers.OverrideWithUint32(s.PUID, other.PUID)
	s.PGID = helpers.OverrideWithUint32(s.PGID, other.PGID)
	s.Timezone = helpers.OverrideWithString(s.Timezone, other.Timezone)
}

func (s *System) setDefaults() {
	const defaultID = 1000
	s.PUID = helpers.DefaultUint32(s.PUID, defaultID)
	s.PGID = helpers.DefaultUint32(s.PGID, defaultID)
}

func (s System) String() string {
	return s.toLinesNode().String()
}

func (s System) toLinesNode() (node *gotree.Node) {
	node = gotree.New("OS Alpine settings:")

	node.Appendf("Process UID: %d", *s.PUID)
	node.Appendf("Process GID: %d", *s.PGID)

	if s.Timezone != "" {
		node.Appendf("Timezone: %s", s.Timezone)
	}

	return node
}
