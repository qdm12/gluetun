package pmtud

type noopLogger struct{}

func (noopLogger) Debug(_ string)            {}
func (noopLogger) Debugf(_ string, _ ...any) {}
func (noopLogger) Warnf(_ string, _ ...any)  {}
