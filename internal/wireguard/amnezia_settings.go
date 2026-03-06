package wireguard

import (
	"fmt"
	"strings"
)

type AmneziaSettings struct {
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

func (s AmneziaSettings) uapiConfig() string {
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
	lines := make([]string, 0, len(uintFields)+len(stringFields))

	for key, val := range uintFields {
		lines = append(lines, fmt.Sprintf("%s=%d", key, val))
	}

	for key, val := range stringFields {
		lines = append(lines, key+"="+val)
	}
	return strings.Join(lines, "\n")
}
