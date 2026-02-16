package netlink

import (
	"fmt"
)

func FamilyToString(family uint8) string {
	switch family {
	case FamilyAll:
		return "all"
	case FamilyV4:
		return "v4"
	case FamilyV6:
		return "v6"
	default:
		return fmt.Sprint(family)
	}
}
