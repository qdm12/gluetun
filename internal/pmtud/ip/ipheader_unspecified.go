//go:build !darwin

package ip

import "encoding/binary"

func putUint16(b []byte, v uint16) {
	binary.BigEndian.PutUint16(b, v)
}
