package nftables

import "github.com/google/nftables"

func IsSupported() bool {
	conn, err := nftables.New()
	if err != nil {
		return false
	}
	_, err = conn.ListTable("filter")
	return err == nil
}
