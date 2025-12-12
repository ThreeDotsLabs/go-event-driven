package log

import (
	"context"
	"log/slog"
	"os"
)

func FromContext(ctx context.Context) *slog.Logger {
	log, ok := ctx.Value(loggerKey).(*slog.Logger)
	if ok {
		return log
	}

	return slog.New(slog.NewTextHandler(os.Stderr, nil))
}

func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
