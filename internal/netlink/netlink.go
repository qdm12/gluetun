package netlink

import "github.com/qdm12/log"

type NetLink struct {
	debugLogger DebugLogger
}

func New(debugLogger DebugLogger) *NetLink {
	return &NetLink{
		debugLogger: debugLogger,
	}
}

func (n *NetLink) PatchLoggerLevel(level log.Level) {
	n.debugLogger.Patch(log.SetLevel(level))
}
