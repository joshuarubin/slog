// Package logfmt implements a "logfmt" format handler.
package logfmt

import (
	"io"
	"os"
	"sync"

	"github.com/go-logfmt/logfmt"
	"jrubin.io/slog"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// Handler implementation.
type Handler struct {
	mu  sync.Mutex
	enc *logfmt.Encoder
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		enc: logfmt.NewEncoder(w),
	}
}

// HandleLog implements slog.Handler.
func (h *Handler) HandleLog(e *slog.Entry) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.enc.EncodeKeyval("time", e.Time)
	h.enc.EncodeKeyval("level", e.Level.String())
	h.enc.EncodeKeyval("message", e.Message)

	for k, v := range e.Fields {
		h.enc.EncodeKeyval(k, v)
	}

	h.enc.EndRecord()

	return nil
}
