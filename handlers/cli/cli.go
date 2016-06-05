// Package cli implements a colored text handler suitable for command-line interfaces.
package cli

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"

	"jrubin.io/slog"
)

// Default handler outputting to stderr.
var Default = New(os.Stderr)

// colors.
const (
	red    = 31
	yellow = 33
	blue   = 34
	gray   = 37
)

// Colors mapping.
var Colors = [...]int{
	slog.DebugLevel: gray,
	slog.InfoLevel:  blue,
	slog.WarnLevel:  yellow,
	slog.ErrorLevel: red,
	slog.FatalLevel: red,
	slog.PanicLevel: red,
}

// Strings mapping.
var Strings = [...]string{
	slog.DebugLevel: "•",
	slog.InfoLevel:  "•",
	slog.WarnLevel:  "•",
	slog.ErrorLevel: "⨯",
	slog.FatalLevel: "⨯",
	slog.PanicLevel: "⨯",
}

// field used for sorting.
type field struct {
	Name  string
	Value interface{}
}

// by sorts projects by call count.
type byName []field

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Handler implementation.
type Handler struct {
	mu      sync.Mutex
	Writer  io.Writer
	Padding int
}

// New handler.
func New(w io.Writer) *Handler {
	return &Handler{
		Writer:  w,
		Padding: 3,
	}
}

var _ slog.Handler = (*Handler)(nil)

// HandleLog implements slog.Handler.
func (h *Handler) HandleLog(e *slog.Entry) error {
	color := Colors[e.Level]
	level := Strings[e.Level]

	var fields []field

	for k, v := range e.Fields {
		fields = append(fields, field{k, v})
	}

	sort.Sort(byName(fields))

	h.mu.Lock()
	defer h.mu.Unlock()

	fmt.Fprintf(h.Writer, "\033[%dm%*s\033[0m %-25s", color, h.Padding+1, level, e.Message)

	for _, f := range fields {
		fmt.Fprintf(h.Writer, " \033[%dm%s\033[0m=%v", color, f.Name, f.Value)
	}

	fmt.Fprintln(h.Writer)

	return nil
}
