package slog_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"jrubin.io/slog"
	"jrubin.io/slog/handlers/discard"
	"jrubin.io/slog/handlers/memory"
	"jrubin.io/slog/handlers/text"
)

type Pet struct {
	Name string
	Age  int
}

func (p *Pet) Fields() slog.Fields {
	return slog.Fields{
		"name": p.Name,
		"age":  p.Age,
	}
}

func TestLogger_printf(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	l.Info("logged in Tobi")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "logged in Tobi")
	assert.Equal(t, e.Level, slog.InfoLevel)
}

func TestFielder(t *testing.T) {
	h := memory.New()
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	pet := &Pet{"Tobi", 3}
	l.WithFields(pet).Info("add pet")

	e := h.Entries[0]
	assert.Equal(t, slog.Fields{"name": "Tobi", "age": 3}, e.Fields)
}

func TestLogger_levels(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	l.Debug("uploading")
	l.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, slog.InfoLevel)
}

func TestLogger_WithFields(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	ctx := l.WithFields(slog.Fields{"file": "sloth.png"})
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, slog.InfoLevel)
	assert.Equal(t, slog.Fields{"file": "sloth.png"}, e.Fields)
}

func TestLogger_WithField(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	ctx := l.WithField("file", "sloth.png").WithField("user", "Tobi")
	ctx.Debug("uploading")
	ctx.Info("upload complete")

	assert.Equal(t, 1, len(h.Entries))

	e := h.Entries[0]
	assert.Equal(t, e.Message, "upload complete")
	assert.Equal(t, e.Level, slog.InfoLevel)
	assert.Equal(t, slog.Fields{"file": "sloth.png", "user": "Tobi"}, e.Fields)
}

func TestLogger_Trace_info(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	func() (err error) {
		defer l.WithField("file", "sloth.png").Trace(slog.InfoLevel, "upload").Stop(&err)
		return nil
	}()

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, slog.InfoLevel)
		assert.Equal(t, slog.Fields{"file": "sloth.png"}, e.Fields)
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, slog.InfoLevel)
		assert.Equal(t, "sloth.png", e.Fields["file"])
		assert.IsType(t, time.Duration(0), e.Fields["duration"])
	}
}

func TestLogger_Trace_error(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	func() (err error) {
		defer l.WithField("file", "sloth.png").Trace(slog.InfoLevel, "upload").Stop(&err)
		return fmt.Errorf("boom")
	}()

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, slog.InfoLevel)
		assert.Equal(t, "sloth.png", e.Fields["file"])
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, slog.ErrorLevel)
		assert.Equal(t, "sloth.png", e.Fields["file"])
		assert.Equal(t, "boom", e.Fields["error"])
		assert.IsType(t, time.Duration(0), e.Fields["duration"])
	}
}

func TestLogger_Trace_nil(t *testing.T) {
	h := memory.New()

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, h)

	func() {
		defer l.WithField("file", "sloth.png").Trace(slog.InfoLevel, "upload").Stop(nil)
	}()

	assert.Equal(t, 2, len(h.Entries))

	{
		e := h.Entries[0]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, slog.InfoLevel)
		assert.Equal(t, slog.Fields{"file": "sloth.png"}, e.Fields)
	}

	{
		e := h.Entries[1]
		assert.Equal(t, e.Message, "upload")
		assert.Equal(t, e.Level, slog.InfoLevel)
		assert.Equal(t, "sloth.png", e.Fields["file"])
		assert.IsType(t, time.Duration(0), e.Fields["duration"])
	}
}

func TestLogger_HandlerFunc(t *testing.T) {
	h := []*slog.Entry{}
	f := func(e *slog.Entry) error {
		h = append(h, e)
		return nil
	}

	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, slog.HandlerFunc(f))

	l.Info("logged in Tobi")

	e := h[0]
	assert.Equal(t, e.Message, "logged in Tobi")
	assert.Equal(t, e.Level, slog.InfoLevel)
}

func BenchmarkLogger_small(b *testing.B) {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, discard.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("login")
	}
}

func BenchmarkLogger_small_text(b *testing.B) {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, text.New(ioutil.Discard))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Info("login")
	}
}

func BenchmarkLogger_medium(b *testing.B) {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, discard.New())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithFields(slog.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).Info("upload")
	}
}

func BenchmarkLogger_large(b *testing.B) {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, discard.New())

	err := fmt.Errorf("boom")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.WithFields(slog.Fields{
			"file": "sloth.png",
			"type": "image/png",
			"size": 1 << 20,
		}).
			WithFields(slog.Fields{
				"some":     "more",
				"data":     "here",
				"whatever": "blah blah",
				"more":     "stuff",
				"context":  "such useful",
				"much":     "fun",
			}).
			WithError(err).Error("upload failed")
	}
}

// Structured logging is supported with fields.
func Example_structured() {
	l := slog.New()
	l.WithField("user", "Tobo").Info("logged in")
}

// Errors are passed to WithError(), populating the "error" field.
func Example_errors() {
	l := slog.New()
	err := errors.New("boom")
	l.WithError(err).Error("upload failed")
}

// Multiple fields can be set, via chaining, or WithFields().
func Example_multipleFields() {
	l := slog.New()
	l.WithFields(slog.Fields{
		"user": "Tobi",
		"file": "sloth.png",
		"type": "image/png",
	}).Info("upload")
}

// Trace can be used to simplify logging of start and completion events,
// for example an upload which may fail.
func Example_trace() {
	var err error
	l := slog.New()
	defer l.Trace(slog.InfoLevel, "upload").Stop(&err)
}
