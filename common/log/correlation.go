package log

import (
	"context"

	"github.com/lithammer/shortuuid/v3"
)

func ContextWithCorrelationID(ctx context.Context, correlationID string) context.Context {
	return context.WithValue(ctx, correlationIDKey, correlationID)
}

func CorrelationIDFromContext(ctx context.Context) string {
	v, ok := ctx.Value(correlationIDKey).(string)
	if ok {
		return v
	}

	FromContext(ctx).Warn("correlation ID not found in context")

	// add "gen_" prefix to distinguish generated correlation IDs from correlation IDs passed by the client
	// it's useful to detect if correlation ID was not passed properly
	return "gen_" + shortuuid.New()
}
