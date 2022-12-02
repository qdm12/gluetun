package netlink

type NetLink struct {
	debugLogger DebugLogger
}

func New(debugLogger DebugLogger) *NetLink {
	return &NetLink{
		debugLogger: debugLogger,
	}
}
