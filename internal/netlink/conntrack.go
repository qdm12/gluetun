package netlink

import (
	"fmt"

	"github.com/ti-mo/conntrack"
)

func (n *NetLink) FlushConntrack() error {
	conn, err := conntrack.Dial(nil)
	if err != nil {
		return fmt.Errorf("dialing conntrack: %w", err)
	}
	defer conn.Close()

	return conn.Flush()
}
