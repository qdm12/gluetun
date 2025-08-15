package netlink

import (
	"fmt"
)

const (
	FamilyAll = 0
	FamilyV4  = 2
	FamilyV6  = 10
)

func FamilyToString(family int) string {
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
