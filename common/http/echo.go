package http

import (
	"errors"
	"net/http"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"
	"github.com/labstack/echo/v4"
)

func NewEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	useMiddlewares(e)
	e.HTTPErrorHandler = HandleError

	return e
}

func HandleError(err error, c echo.Context) {
	log.FromContext(c.Request().Context()).With("error", err).Error("HTTP error")

	httpCode := http.StatusInternalServerError
	msg := any("Internal server error")

	httpErr := &echo.HTTPError{}
	if errors.As(err, &httpErr) {
		httpCode = httpErr.Code
		msg = httpErr.Message
	}

	jsonErr := c.JSON(
		httpCode,
		map[string]any{
			"error": msg,
		},
	)
	if jsonErr != nil {
		panic(err)
	}
}
