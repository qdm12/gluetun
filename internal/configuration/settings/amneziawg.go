package settings

import (
	"errors"
	"fmt"

	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

type AmneziaWg struct {
	JunkPacketCount uint16 `json:"junk_packet_count"`
	JunkPacketMin   uint16 `json:"junk_packet_min"`
	JunkPacketMax   uint16 `json:"junk_packet_max"`
	PaddingS1       uint16 `json:"padding_s1"`
	PaddingS2       uint16 `json:"padding_s2"`
	PaddingS3       uint16 `json:"padding_s3"`
	PaddingS4       uint16 `json:"padding_s4"`
	HeaderH1        string `json:"header_h1"`
	HeaderH2        string `json:"header_h2"`
	HeaderH3        string `json:"header_h3"`
	HeaderH4        string `json:"header_h4"`
	InitPacketI1    string `json:"init_packet_i1"`
	InitPacketI2    string `json:"init_packet_i2"`
	InitPacketI3    string `json:"init_packet_i3"`
	InitPacketI4    string `json:"init_packet_i4"`
	InitPacketI5    string `json:"init_packet_i5"`
}

func (s AmneziaWg) copy() (copied AmneziaWg) {
	copied.JunkPacketCount = s.JunkPacketCount
	copied.JunkPacketMin = s.JunkPacketMin
	copied.JunkPacketMax = s.JunkPacketMax
	copied.PaddingS1 = s.PaddingS1
	copied.PaddingS2 = s.PaddingS2
	copied.PaddingS3 = s.PaddingS3
	copied.PaddingS4 = s.PaddingS4
	copied.HeaderH1 = s.HeaderH1
	copied.HeaderH2 = s.HeaderH2
	copied.HeaderH3 = s.HeaderH3
	copied.HeaderH4 = s.HeaderH4
	copied.InitPacketI1 = s.InitPacketI1
	copied.InitPacketI2 = s.InitPacketI2
	copied.InitPacketI3 = s.InitPacketI3
	copied.InitPacketI4 = s.InitPacketI4
	copied.InitPacketI5 = s.InitPacketI5
	return copied
}

//nolint:dupl
func (s *AmneziaWg) overrideWith(other AmneziaWg) {
	s.JunkPacketCount = gosettings.OverrideWithComparable(s.JunkPacketCount, other.JunkPacketCount)
	s.JunkPacketMin = gosettings.OverrideWithComparable(s.JunkPacketMin, other.JunkPacketMin)
	s.JunkPacketMax = gosettings.OverrideWithComparable(s.JunkPacketMax, other.JunkPacketMax)
	s.PaddingS1 = gosettings.OverrideWithComparable(s.PaddingS1, other.PaddingS1)
	s.PaddingS2 = gosettings.OverrideWithComparable(s.PaddingS2, other.PaddingS2)
	s.PaddingS3 = gosettings.OverrideWithComparable(s.PaddingS3, other.PaddingS3)
	s.PaddingS4 = gosettings.OverrideWithComparable(s.PaddingS4, other.PaddingS4)
	s.HeaderH1 = gosettings.OverrideWithComparable(s.HeaderH1, other.HeaderH1)
	s.HeaderH2 = gosettings.OverrideWithComparable(s.HeaderH2, other.HeaderH2)
	s.HeaderH3 = gosettings.OverrideWithComparable(s.HeaderH3, other.HeaderH3)
	s.HeaderH4 = gosettings.OverrideWithComparable(s.HeaderH4, other.HeaderH4)
	s.InitPacketI1 = gosettings.OverrideWithComparable(s.InitPacketI1, other.InitPacketI1)
	s.InitPacketI2 = gosettings.OverrideWithComparable(s.InitPacketI2, other.InitPacketI2)
	s.InitPacketI3 = gosettings.OverrideWithComparable(s.InitPacketI3, other.InitPacketI3)
	s.InitPacketI4 = gosettings.OverrideWithComparable(s.InitPacketI4, other.InitPacketI4)
	s.InitPacketI5 = gosettings.OverrideWithComparable(s.InitPacketI5, other.InitPacketI5)
}

func (s AmneziaWg) toLinesNode() (node *gotree.Node) {
	node = gotree.New("Amneziawg parameters:")

	uintFields := []struct {
		key string
		val uint16
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
		if f.val != 0 {
			node.Appendf("%s: %d", f.key, f.val)
		}
	}
	stringFields := []struct {
		key string
		val string
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
		if f.val != "" {
			node.Appendf("%s: %s", f.key, f.val)
		}
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
	case s.JunkPacketMax != 0 && s.JunkPacketMin > s.JunkPacketMax:
		return fmt.Errorf("%w: jmin=%d and jmax=%d",
			ErrJunkPacketBounds, s.JunkPacketMin, s.JunkPacketMax)
	case s.JunkPacketCount == 0 && (s.JunkPacketMin != 0 || s.JunkPacketMax != 0):
		return fmt.Errorf("%w: jc=%d and jmin=%d and jmax=%d",
			ErrJunkPacketCountNotSet, s.JunkPacketCount, s.JunkPacketMin, s.JunkPacketMax)
	case s.JunkPacketCount != 0 && (s.JunkPacketMin == 0 || s.JunkPacketMax == 0):
		return fmt.Errorf("%w: jc=%d and jmin=%d and jmax=%d",
			ErrJunkPacketMinMaxNotSet, s.JunkPacketCount, s.JunkPacketMin, s.JunkPacketMax)
	}

	return nil
}
