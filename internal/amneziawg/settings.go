package amneziawg

import (
	"errors"
	"fmt"
	"strings"
)

type Settings struct {
	JunkPacketCount uint16
	JunkPacketMin   uint16
	JunkPacketMax   uint16
	PaddingS1       uint16
	PaddingS2       uint16
	PaddingS3       uint16
	PaddingS4       uint16
	HeaderH1        string
	HeaderH2        string
	HeaderH3        string
	HeaderH4        string
	InitPacketI1    string
	InitPacketI2    string
	InitPacketI3    string
	InitPacketI4    string
	InitPacketI5    string
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
