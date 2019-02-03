package main

import (
	"errors"
	"os"
	"time"

	"github.com/joshuarubin/slog"
	"github.com/joshuarubin/slog/handlers/logfmt"
)

func main() {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, logfmt.New(os.Stderr))

	ctx := l.WithFields(slog.Fields{
		"file": "something.png",
		"type": "image/png",
		"user": "tobi",
	})

	for range time.Tick(time.Millisecond * 200) {
		ctx.Info("upload")
		ctx.Info("upload complete")
		ctx.Warn("upload retry")
		ctx.WithError(errors.New("unauthorized")).Error("upload failed")
	}
}
