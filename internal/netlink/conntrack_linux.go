package netlink

import (
	"fmt"

	"github.com/mdlayher/netlink"
	"github.com/ti-mo/netfilter"
)

func (n *NetLink) FlushConntrack() error {
	conn, err := netfilter.Dial(nil)
	if err != nil {
		return fmt.Errorf("dialing netfilter: %w", err)
	}
	defer conn.Close()

	families := [...]netfilter.ProtoFamily{netfilter.ProtoIPv4, netfilter.ProtoIPv6}
	for _, family := range families {
		const IPCtnlMsgCtDelete = 2
		request, err := netfilter.MarshalNetlink(
			netfilter.Header{
				SubsystemID: netfilter.NFSubsysCTNetlink,
				MessageType: netfilter.MessageType(IPCtnlMsgCtDelete),
				Family:      family,
				Flags:       netlink.Request | netlink.Acknowledge,
			},
			nil)
		if err != nil {
			return fmt.Errorf("encoding netlink request: %w", err)
		}

		_, err = conn.Query(request)
		if err != nil {
			return fmt.Errorf("querying netlink request: %w", err)
		}
	}
	return nil
}
