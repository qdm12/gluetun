package firewall

import (
	"errors"
	"fmt"
)

type tcpFlags struct {
	mask       []tcpFlag
	comparison []tcpFlag
}

type tcpFlag uint8

const (
	tcpFlagFIN tcpFlag = 1 << iota
	tcpFlagSYN
	tcpFlagRST
	tcpFlagPSH
	tcpFlagACK
	tcpFlagURG
	tcpFlagECE
	tcpFlagCWR
)

func (f tcpFlag) String() string {
	switch f {
	case tcpFlagFIN:
		return "FIN"
	case tcpFlagSYN:
		return "SYN"
	case tcpFlagRST:
		return "RST"
	case tcpFlagPSH:
		return "PSH"
	case tcpFlagACK:
		return "ACK"
	case tcpFlagURG:
		return "URG"
	case tcpFlagECE:
		return "ECE"
	case tcpFlagCWR:
		return "CWR"
	default:
		panic(fmt.Sprintf("%s: %d", errTCPFlagUnknown, f))
	}
}

var errTCPFlagUnknown = errors.New("unknown TCP flag")

func parseTCPFlag(s string) (tcpFlag, error) {
	allFlags := []tcpFlag{
		tcpFlagFIN, tcpFlagSYN, tcpFlagRST, tcpFlagPSH,
		tcpFlagACK, tcpFlagURG, tcpFlagECE, tcpFlagCWR,
	}
	for _, flag := range allFlags {
		if s == fmt.Sprintf("%#02x", uint8(flag)) || s == flag.String() {
			return flag, nil
		}
	}
	return 0, fmt.Errorf("%w: %s", errTCPFlagUnknown, s)
}
