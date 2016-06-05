package main

import (
	"errors"
	"time"

	"jrubin.io/slog"
	"jrubin.io/slog/handlers/auto"
)

func main() {
	l := slog.New()
	l.RegisterHandler(slog.InfoLevel, auto.Default)

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
		ctx.WithField("file", "img.png").Error("failed to upload")
	}
}
