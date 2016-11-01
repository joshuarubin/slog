package slog

import (
	"bytes"
	"io"
	"runtime"
	"sync"
)

type syncWriter struct {
	PrintFunc func(string)
	buf       bytes.Buffer
	mu        sync.Mutex
}

func (w *syncWriter) printLines() {
	for {
		i := bytes.IndexByte(w.buf.Bytes(), '\n')
		if i < 0 {
			break
		}

		data := w.buf.Next(i + 1)

		// strip trailing "\r\n"
		for _, b := range []byte("\n\r") { // yes, I know "\n\r" is backwards
			if len(data) > 0 && data[len(data)-1] == b {
				data = data[:len(data)-1]
			}
		}

		w.PrintFunc(string(data))
	}
}

func (w *syncWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	n, err := w.buf.Write(p)
	w.printLines()
	w.mu.Unlock()

	return n, err
}

func (w *syncWriter) Close() error {
	w.mu.Lock()
	w.printLines()
	w.PrintFunc(w.buf.String())
	w.buf.Reset()
	w.mu.Unlock()
	return nil
}

// Writer returns an io.WriteCloser where each line written to that writer will
// be printed using the handlers for the given Level. It is the caller's
// responsibility to close it.
func (logger *Logger) Writer(level Level) io.WriteCloser {
	if level < PanicLevel {
		level = PanicLevel
	}

	if level > DebugLevel {
		level = DebugLevel
	}

	var printFunc func(msg string)
	switch level {
	case DebugLevel:
		printFunc = logger.Debug
	case InfoLevel:
		printFunc = logger.Info
	case WarnLevel:
		printFunc = logger.Warn
	case ErrorLevel:
		printFunc = logger.Error
	case FatalLevel:
		printFunc = logger.Fatal
	case PanicLevel:
		printFunc = logger.Panic
	}

	w := &syncWriter{
		PrintFunc: printFunc,
	}

	runtime.SetFinalizer(w, writerFinalizer)

	return w
}

func writerFinalizer(writer *syncWriter) {
	_ = writer.Close()
}
