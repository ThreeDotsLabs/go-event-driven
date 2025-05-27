package http

import (
	"log/slog"
	"unicode/utf8"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lithammer/shortuuid/v3"
)

func useMiddlewares(e *echo.Echo) {
	e.Use(
		middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: func() string {
				return shortuuid.New()
			},
		}),
		middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)

			logger := log.FromContext(c.Request().Context()).With(
				"request_id", reqID,
				"request_body: ", string(reqBody),
			)

			if utf8.ValidString(string(resBody)) {
				logger = logger.With("response_body: ", string(resBody))
			} else {
				logger = logger.With("response_body: ", "<binary data>")
			}

			logger.Info("Request/response")
		}),
		middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogURI:       true,
			LogRequestID: true,
			LogStatus:    true,
			LogMethod:    true,
			LogLatency:   true,
			LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
				log.FromContext(c.Request().Context()).With(
					"URI", values.URI,
					"request_id", values.RequestID,
					"status", values.Status,
					"method", values.Method,
					"duration", values.Latency.String(),
					"error", values.Error,
				).Info("Request done")

				return nil
			},
		}),
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				req := c.Request()
				ctx := req.Context()

				reqCorrelationID := req.Header.Get(log.CorrelationIDHttpHeader)
				if reqCorrelationID == "" {
					reqCorrelationID = shortuuid.New()
				}

				ctx = log.ToContext(ctx, slog.With("correlation_id", reqCorrelationID))
				ctx = log.ContextWithCorrelationID(ctx, reqCorrelationID)

				c.SetRequest(req.WithContext(ctx))
				c.Response().Header().Set(log.CorrelationIDHttpHeader, reqCorrelationID)

				return next(c)
			}
		},
	)
}
