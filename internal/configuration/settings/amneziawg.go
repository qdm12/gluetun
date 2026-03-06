package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gosettings/reader"
	"github.com/qdm12/gotree"
)

type AmneziaWg struct {
	JunkPacketCount *uint16 `json:"junk_packet_count"`
	JunkPacketMin   *uint16 `json:"junk_packet_min"`
	JunkPacketMax   *uint16 `json:"junk_packet_max"`
	PaddingS1       *uint16 `json:"padding_s1"`
	PaddingS2       *uint16 `json:"padding_s2"`
	PaddingS3       *uint16 `json:"padding_s3"`
	PaddingS4       *uint16 `json:"padding_s4"`
	HeaderH1        *string `json:"header_h1"`
	HeaderH2        *string `json:"header_h2"`
	HeaderH3        *string `json:"header_h3"`
	HeaderH4        *string `json:"header_h4"`
	InitPacketI1    *string `json:"init_packet_i1"`
	InitPacketI2    *string `json:"init_packet_i2"`
	InitPacketI3    *string `json:"init_packet_i3"`
	InitPacketI4    *string `json:"init_packet_i4"`
	InitPacketI5    *string `json:"init_packet_i5"`
}

func (s *AmneziaWg) read(r *reader.Reader) error {
	uint16Fields := map[string]**uint16{
		"AMNEZIAWG_JC":   &s.JunkPacketCount,
		"AMNEZIAWG_JMIN": &s.JunkPacketMin,
		"AMNEZIAWG_JMAX": &s.JunkPacketMax,
		"AMNEZIAWG_S1":   &s.PaddingS1,
		"AMNEZIAWG_S2":   &s.PaddingS2,
		"AMNEZIAWG_S3":   &s.PaddingS3,
		"AMNEZIAWG_S4":   &s.PaddingS4,
	}
	for key, dst := range uint16Fields {
		v, err := r.Uint16Ptr(key)
		if err != nil {
			return err
		}
		*dst = v
	}
	stringFields := map[string]**string{
		"AMNEZIAWG_H1": &s.HeaderH1,
		"AMNEZIAWG_H2": &s.HeaderH2,
		"AMNEZIAWG_H3": &s.HeaderH3,
		"AMNEZIAWG_H4": &s.HeaderH4,
		"AMNEZIAWG_I1": &s.InitPacketI1,
		"AMNEZIAWG_I2": &s.InitPacketI2,
		"AMNEZIAWG_I3": &s.InitPacketI3,
		"AMNEZIAWG_I4": &s.InitPacketI4,
		"AMNEZIAWG_I5": &s.InitPacketI5,
	}
	opt := reader.ForceLowercase(false)
	for key, dst := range stringFields {
		*dst = r.Get(key, opt) // *string (nil если не задано)
	}
	return nil
}

func (s AmneziaWg) copy() (copied AmneziaWg) {
	return AmneziaWg{
		JunkPacketCount: gosettings.CopyPointer(s.JunkPacketCount),
		JunkPacketMin:   gosettings.CopyPointer(s.JunkPacketMin),
		JunkPacketMax:   gosettings.CopyPointer(s.JunkPacketMax),
		PaddingS1:       gosettings.CopyPointer(s.PaddingS1),
		PaddingS2:       gosettings.CopyPointer(s.PaddingS2),
		PaddingS3:       gosettings.CopyPointer(s.PaddingS3),
		PaddingS4:       gosettings.CopyPointer(s.PaddingS4),
		HeaderH1:        gosettings.CopyPointer(s.HeaderH1),
		HeaderH2:        gosettings.CopyPointer(s.HeaderH2),
		HeaderH3:        gosettings.CopyPointer(s.HeaderH3),
		HeaderH4:        gosettings.CopyPointer(s.HeaderH4),
		InitPacketI1:    gosettings.CopyPointer(s.InitPacketI1),
		InitPacketI2:    gosettings.CopyPointer(s.InitPacketI2),
		InitPacketI3:    gosettings.CopyPointer(s.InitPacketI3),
		InitPacketI4:    gosettings.CopyPointer(s.InitPacketI4),
		InitPacketI5:    gosettings.CopyPointer(s.InitPacketI5),
	}
}

//nolint:dupl
func (s *AmneziaWg) overrideWith(other AmneziaWg) {
	s.JunkPacketCount = gosettings.OverrideWithPointer(s.JunkPacketCount, other.JunkPacketCount)
	s.JunkPacketMin = gosettings.OverrideWithPointer(s.JunkPacketMin, other.JunkPacketMin)
	s.JunkPacketMax = gosettings.OverrideWithPointer(s.JunkPacketMax, other.JunkPacketMax)
	s.PaddingS1 = gosettings.OverrideWithPointer(s.PaddingS1, other.PaddingS1)
	s.PaddingS2 = gosettings.OverrideWithPointer(s.PaddingS2, other.PaddingS2)
	s.PaddingS3 = gosettings.OverrideWithPointer(s.PaddingS3, other.PaddingS3)
	s.PaddingS4 = gosettings.OverrideWithPointer(s.PaddingS4, other.PaddingS4)
	s.HeaderH1 = gosettings.OverrideWithPointer(s.HeaderH1, other.HeaderH1)
	s.HeaderH2 = gosettings.OverrideWithPointer(s.HeaderH2, other.HeaderH2)
	s.HeaderH3 = gosettings.OverrideWithPointer(s.HeaderH3, other.HeaderH3)
	s.HeaderH4 = gosettings.OverrideWithPointer(s.HeaderH4, other.HeaderH4)
	s.InitPacketI1 = gosettings.OverrideWithPointer(s.InitPacketI1, other.InitPacketI1)
	s.InitPacketI2 = gosettings.OverrideWithPointer(s.InitPacketI2, other.InitPacketI2)
	s.InitPacketI3 = gosettings.OverrideWithPointer(s.InitPacketI3, other.InitPacketI3)
	s.InitPacketI4 = gosettings.OverrideWithPointer(s.InitPacketI4, other.InitPacketI4)
	s.InitPacketI5 = gosettings.OverrideWithPointer(s.InitPacketI5, other.InitPacketI5)
}

func (s *AmneziaWg) setDefaults() {
	s.JunkPacketCount = gosettings.DefaultPointer(s.JunkPacketCount, 0)
	s.JunkPacketMin = gosettings.DefaultPointer(s.JunkPacketMin, 0)
	s.JunkPacketMax = gosettings.DefaultPointer(s.JunkPacketMax, 0)
	s.PaddingS1 = gosettings.DefaultPointer(s.PaddingS1, 0)
	s.PaddingS2 = gosettings.DefaultPointer(s.PaddingS2, 0)
	s.PaddingS3 = gosettings.DefaultPointer(s.PaddingS3, 0)
	s.PaddingS4 = gosettings.DefaultPointer(s.PaddingS4, 0)
	s.HeaderH1 = gosettings.DefaultPointer(s.HeaderH1, "")
	s.HeaderH2 = gosettings.DefaultPointer(s.HeaderH2, "")
	s.HeaderH3 = gosettings.DefaultPointer(s.HeaderH3, "")
	s.HeaderH4 = gosettings.DefaultPointer(s.HeaderH4, "")
	s.InitPacketI1 = gosettings.DefaultPointer(s.InitPacketI1, "")
	s.InitPacketI2 = gosettings.DefaultPointer(s.InitPacketI2, "")
	s.InitPacketI3 = gosettings.DefaultPointer(s.InitPacketI3, "")
	s.InitPacketI4 = gosettings.DefaultPointer(s.InitPacketI4, "")
	s.InitPacketI5 = gosettings.DefaultPointer(s.InitPacketI5, "")
}

func (s AmneziaWg) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Amneziawg parameters:")

	uintFields := []struct {
		key string
		val *uint16
	}{
		{"jc", s.JunkPacketCount},
		{"jmin", s.JunkPacketMin},
		{"jmax", s.JunkPacketMax},
		{"s1", s.PaddingS1},
		{"s2", s.PaddingS2},
		{"s3", s.PaddingS3},
		{"s4", s.PaddingS4},
	}
	for _, f := range uintFields {
		node.Appendf("%s: %d", f.key, *f.val)
	}

	stringFields := []struct {
		key string
		val *string
	}{
		{"h1", s.HeaderH1},
		{"h2", s.HeaderH2},
		{"h3", s.HeaderH3},
		{"h4", s.HeaderH4},
		{"i1", s.InitPacketI1},
		{"i2", s.InitPacketI2},
		{"i3", s.InitPacketI3},
		{"i4", s.InitPacketI4},
		{"i5", s.InitPacketI5},
	}
	for _, f := range stringFields {
		node.Appendf("%s: %s", f.key, *f.val)
	}

	return node
}

var (
	ErrJunkPacketBounds       = errors.New("junk packet minimum must be lower than or equal to maximum")
	ErrJunkPacketMinMaxNotSet = errors.New("junk packet min and max must be set when junk packet count is set")
	ErrJunkPacketCountNotSet  = errors.New("junk packet count must be set when junk packet min or max is set")
)

func (s AmneziaWg) validate() error {
	switch {
	case *s.JunkPacketMax != 0 && *s.JunkPacketMin > *s.JunkPacketMax:
		return fmt.Errorf("%w: jmin=%d and jmax=%d",
			ErrJunkPacketBounds, s.JunkPacketMin, s.JunkPacketMax)
	case *s.JunkPacketCount == 0 && (*s.JunkPacketMin != 0 || *s.JunkPacketMax != 0):
		return fmt.Errorf("%w: jc=%d and jmin=%d and jmax=%d",
			ErrJunkPacketCountNotSet, s.JunkPacketCount, *s.JunkPacketMin, *s.JunkPacketMax)
	case *s.JunkPacketCount != 0 && (*s.JunkPacketMin == 0 || *s.JunkPacketMax == 0):
		return fmt.Errorf("%w: jc=%d and jmin=%d and jmax=%d",
			ErrJunkPacketMinMaxNotSet, s.JunkPacketCount, s.JunkPacketMin, s.JunkPacketMax)
	}

	return nil
}
