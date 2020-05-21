package calculations

import (
	"context"
)

type Logger interface {
	Infof(message string, args ...interface{})
	Warnf(message string, args ...interface{})
	Debugf(message string, args ...interface{})
	Errorf(message string, args ...interface{})
}

type noopLogger struct{}

func (n *noopLogger) Infof(string, ...interface{}) {
	// do nothing
}
func (n *noopLogger) Warnf(string, ...interface{}) {
	// do nothing
}
func (n *noopLogger) Debugf(string, ...interface{}) {
	// do nothing
}
func (n *noopLogger) Errorf(string, ...interface{}) {
	// do nothing
}

func newNoopLogger() *noopLogger {
	return &noopLogger{}
}

type runtime struct {
	Ctx    context.Context
	Logger Logger
}
