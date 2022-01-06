package helpers

func BoolPtrToYesNo(b *bool) string {
	if *b {
		return "yes"
	}
	return "no"
}

func TCPPtrToString(tcp *bool) string {
	if *tcp {
		return "TCP"
	}
	return "UDP"
}
