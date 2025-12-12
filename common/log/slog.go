package log

import (
	"log/slog"
	"os"

	"github.com/ThreeDotsLabs/humanslog"
)

func Init(level slog.Level) {
	opts := &humanslog.Options{
		HandlerOptions: &slog.HandlerOptions{
			Level: level,
		},
		TimeFormat: "[15:04:05.000]",
	}

	logger := slog.New(humanslog.NewHandler(os.Stderr, opts))

	slog.SetDefault(logger)
}
