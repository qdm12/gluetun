package ip

import (
	"encoding/binary"
)

func putUint16(b []byte, v uint16) {
	binary.NativeEndian.PutUint16(b, v)
}
