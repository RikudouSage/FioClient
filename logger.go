package fioclient

// NoopGooseLogger discards all log messages emitted by goose migrations.
type NoopGooseLogger struct {
}

// Fatalf discards a formatted fatal log message.
func (receiver *NoopGooseLogger) Fatalf(format string, v ...any) {
}

// Printf discards a formatted log message.
func (receiver *NoopGooseLogger) Printf(format string, v ...any) {
}

// NewNoopGooseLogger creates a goose logger that discards all messages.
func NewNoopGooseLogger() *NoopGooseLogger {
	return &NoopGooseLogger{}
}
