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
	ErrJunkPacketBounds = errors.New("junk packet minimum must be lower or equal than maximum")
	ErrJunkPacketSize   = errors.New("junk packet min and max must be set when junk packet count is set")
	ErrJunkPacketCount  = errors.New("junk packet count must be set when junk packet min or max is set")
)

func (s Settings) Validate() error {
	if s.JunkPacketMin > s.JunkPacketMax && s.JunkPacketMax != 0 {
		return fmt.Errorf(
			"%w: jmin=%d and jmax=%d",
			ErrJunkPacketBounds,
			s.JunkPacketMin,
			s.JunkPacketMax,
		)
	}

	hasJunkSize := s.JunkPacketMin != 0 || s.JunkPacketMax != 0
	if s.JunkPacketCount == 0 && hasJunkSize {
		return fmt.Errorf(
			"%w: jc=%d and jmin=%d and jmax=%d",
			ErrJunkPacketCount,
			s.JunkPacketCount,
			s.JunkPacketMin,
			s.JunkPacketMax,
		)
	}

	if s.JunkPacketCount != 0 && (s.JunkPacketMin == 0 || s.JunkPacketMax == 0) {
		return fmt.Errorf(
			"%w: jc=%d and jmin=%d and jmax=%d",
			ErrJunkPacketSize,
			s.JunkPacketCount,
			s.JunkPacketMin,
			s.JunkPacketMax,
		)
	}

	return nil
}

func (s Settings) UAPIConfig() string {
	var lines []string
	if s.JunkPacketCount != 0 {
		lines = append(lines, fmt.Sprintf("jc=%d", s.JunkPacketCount))
	}
	if s.JunkPacketMin != 0 {
		lines = append(lines, fmt.Sprintf("jmin=%d", s.JunkPacketMin))
	}
	if s.JunkPacketMax != 0 {
		lines = append(lines, fmt.Sprintf("jmax=%d", s.JunkPacketMax))
	}
	if s.PaddingS1 != 0 {
		lines = append(lines, fmt.Sprintf("s1=%d", s.PaddingS1))
	}
	if s.PaddingS2 != 0 {
		lines = append(lines, fmt.Sprintf("s2=%d", s.PaddingS2))
	}
	if s.PaddingS3 != 0 {
		lines = append(lines, fmt.Sprintf("s3=%d", s.PaddingS3))
	}
	if s.PaddingS4 != 0 {
		lines = append(lines, fmt.Sprintf("s4=%d", s.PaddingS4))
	}
	if s.HeaderH1 != "" {
		lines = append(lines, "h1="+s.HeaderH1)
	}
	if s.HeaderH2 != "" {
		lines = append(lines, "h2="+s.HeaderH2)
	}
	if s.HeaderH3 != "" {
		lines = append(lines, "h3="+s.HeaderH3)
	}
	if s.HeaderH4 != "" {
		lines = append(lines, "h4="+s.HeaderH4)
	}
	if s.InitPacketI1 != "" {
		lines = append(lines, "i1="+s.InitPacketI1)
	}
	if s.InitPacketI2 != "" {
		lines = append(lines, "i2="+s.InitPacketI2)
	}
	if s.InitPacketI3 != "" {
		lines = append(lines, "i3="+s.InitPacketI3)
	}
	if s.InitPacketI4 != "" {
		lines = append(lines, "i4="+s.InitPacketI4)
	}
	if s.InitPacketI5 != "" {
		lines = append(lines, "i5="+s.InitPacketI5)
	}

	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}
