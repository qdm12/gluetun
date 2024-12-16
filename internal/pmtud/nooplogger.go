package pmtud

type noopLogger struct{}

func (noopLogger) Debug(_ string, _ ...any) {}
