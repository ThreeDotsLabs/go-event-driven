package log

import (
	"context"

	"github.com/sirupsen/logrus"
)

func FromContext(ctx context.Context) *logrus.Entry {
	log, ok := ctx.Value(loggerKey).(*logrus.Entry)
	if ok {
		return log
	}

	return logrus.NewEntry(logrus.StandardLogger())
}

func ToContext(ctx context.Context, log *logrus.Entry) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}
