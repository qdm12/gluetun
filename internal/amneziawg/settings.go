package amneziawg

import (
	"errors"
	"fmt"
	"strings"
)

type Settings struct {
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

func (s Settings) IsZero() bool {
	return s == Settings{}
}

var (
	ErrJunkPacketBounds       = errors.New("junk packet minimum must be lower than or equal to maximum")
	ErrJunkPacketMinMaxNotSet = errors.New("junk packet min and max must be set when junk packet count is set")
	ErrJunkPacketCountNotSet  = errors.New("junk packet count must be set when junk packet min or max is set")
)

func (s Settings) Validate() error {
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

func (s Settings) UAPIConfig() string {
	var lines []string
	uintFields := map[string]uint16{
		"jc":   s.JunkPacketCount,
		"jmin": s.JunkPacketMin,
		"jmax": s.JunkPacketMax,
		"s1":   s.PaddingS1,
		"s2":   s.PaddingS2,
		"s3":   s.PaddingS3,
		"s4":   s.PaddingS4,
	}
	stringFields := map[string]string{
		"h1": s.HeaderH1,
		"h2": s.HeaderH2,
		"h3": s.HeaderH3,
		"h4": s.HeaderH4,
		"i1": s.InitPacketI1,
		"i2": s.InitPacketI2,
		"i3": s.InitPacketI3,
		"i4": s.InitPacketI4,
		"i5": s.InitPacketI5,
	}

	for key, val := range uintFields {
		if val != 0 {
			lines = append(lines, fmt.Sprintf("%s=%d", key, val))
		}
	}

	for key, val := range stringFields {
		if val != "" {
			lines = append(lines, key+"="+val)
		}
	}
	return strings.Join(lines, "\n")
}
