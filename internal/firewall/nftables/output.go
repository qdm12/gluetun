package nftables

import (
	"context"
	"fmt"

	"github.com/google/nftables"
	"github.com/google/nftables/expr"
)

func (f *Firewall) AcceptIpv6MulticastOutput(_ context.Context, intf string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	conn, err := nftables.New()
	if err != nil {
		return fmt.Errorf("creating nftables connection: %w", err)
	}

	table, _, _, outputChain := setupFilterWithBaseChains(conn)

	const maxExprsLen = 6
	exprs := make([]expr.Any, 0, maxExprsLen)

	if intf != "" && intf != "*" {
		exprs = append(exprs,
			&expr.Meta{Key: expr.MetaKeyOIFNAME, Register: 1},
			&expr.Cmp{Op: expr.CmpOpEq, Register: 1, Data: []byte(intf + "\x00")},
		)
	}

	// ff02::1:ff00:0/104 mask is 13 bytes of 0xff
	mask := []byte{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00,
	} //nolint:mnd
	addr := []byte{
		0xff, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01, 0xff, 0x00, 0x00, 0x00,
	} //nolint:mnd

	exprs = append(exprs,
		&expr.Payload{
			DestRegister: 1,
			Base:         expr.PayloadBaseNetworkHeader,
			Offset:       24, // IPv6 Destination Address offset //nolint:mnd
			Len:          16, //nolint:mnd
		},
		&expr.Bitwise{
			SourceRegister: 1,
			DestRegister:   1,
			Len:            16, //nolint:mnd
			Mask:           mask,
			Xor:            []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, //nolint:mnd
		},
		&expr.Cmp{
			Op:       expr.CmpOpEq,
			Register: 1,
			Data:     addr,
		},
		&expr.Verdict{Kind: expr.VerdictAccept},
	)

	rule := &nftables.Rule{
		Table: table,
		Chain: outputChain,
		Exprs: exprs,
	}

	conn.AddRule(rule)

	err = conn.Flush()
	if err != nil {
		return fmt.Errorf("flushing: %w", err)
	}

	return nil
}
