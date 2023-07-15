package utils

import (
	"context"
)

type contextKey string

const (
	// ContextDebug is a key for the debug flag in a context
	ContextDebug = contextKey("debug")
)

// DebugFromContext get the debug flag from a context
func DebugFromContext(ctx context.Context) bool {
	var (
		debug bool = false
	)
	ctxDebug := ctx.Value(ContextDebug)
	if ctxDebug != nil {
		debug = ctxDebug.(bool)
	}
	return debug
}
