package helpers

func TCPPtrToString(tcp *bool) string {
	if *tcp {
		return "TCP"
	}
	return "UDP"
}
