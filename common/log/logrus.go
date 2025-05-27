package log

import (
	"log/slog"
	"os"
)

func Init(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	logger := slog.New(slog.NewTextHandler(os.Stderr, opts))

	slog.SetDefault(logger)
}
