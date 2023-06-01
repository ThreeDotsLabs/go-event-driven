package log

type ctxKey int

const (
	loggerKey        ctxKey = iota
	correlationIDKey ctxKey = iota
)
