// Package discard implements a no-op handler useful for benchmarks and tests.
package discard

import "jrubin.io/slog"

// Default handler.
var Default = New()

// Handler implementation.
type Handler struct{}

// New handler.
func New() *Handler {
	return &Handler{}
}

// HandleLog implements slog.Handler.
func (h *Handler) HandleLog(e *slog.Entry) error {
	return nil
}
