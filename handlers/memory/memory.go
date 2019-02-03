// Package memory implements an in-memory handler useful for testing, as the
// entries can be accessed after writes.
package memory

import (
	"sync"

	"github.com/joshuarubin/slog"
)

// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Entries []*slog.Entry
}

// New handler.
func New() *Handler {
	return &Handler{}
}

// HandleLog implements slog.Handler.
func (h *Handler) HandleLog(e *slog.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.Entries = append(h.Entries, e)
	return nil
}
