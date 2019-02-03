package main

import (
	"os"
	"time"

	"github.com/joshuarubin/slog"
	"github.com/joshuarubin/slog/handlers/text"
)

func work(ctx slog.Interface) (err error) {
	path := "README.md"
	defer ctx.WithField("path", path).Trace(slog.InfoLevel, "opening").Stop(&err)
	_, err = os.Open(path)
	return
}

func main() {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, text.New(os.Stderr))

	ctx := l.WithFields(slog.Fields{
		"app": "myapp",
		"env": "prod",
	})

	for range time.Tick(time.Second) {
		_ = work(ctx)
	}
}
