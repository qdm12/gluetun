package tcp

func stripIPv4Header(reply []byte) (result []byte, ok bool) {
	return reply, true
}
