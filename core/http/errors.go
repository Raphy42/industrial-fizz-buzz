package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/Raphy42/industrial-fizz-buzz/core/config"
	"github.com/Raphy42/industrial-fizz-buzz/core/errors"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

func errorBody(value any) map[string]any {
	return map[string]any{
		"error": value,
	}
}

func developmentErrorMiddleware(err error, c echo.Context) {
	var body any
	status := http.StatusInternalServerError
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	switch v := err.(type) {
	case *errors.Error:
		body = v
		status = v.HttpCode
	case error:
		body = errorBody(err)
	default:
		body = errorBody("internal server error")
	}
	log.Error("handler error",
		zap.String("request.path", c.Path()),
		zap.String("request.method", c.Request().Method),
		zap.String("request.uri", c.Request().RequestURI),
		zap.Error(err),
	)
	if err = c.JSON(status, body); err != nil {
		log.Fatal("unrecoverable error while serializing error response", zap.Error(err))
	}
}

func prodErrorMiddleware(err error, c echo.Context) {
	var body any
	status := http.StatusInternalServerError
	ctx := c.Request().Context()
	log := logger.FromContext(ctx)

	switch v := err.(type) {
	case *errors.Error:
		body = errorBody(v.Message)
	default:
		body = errorBody("internal server error")
	}
	log.Error("handler error",
		zap.String("request.path", c.Path()),
		zap.String("request.method", c.Request().Method),
		zap.String("request.uri", c.Request().RequestURI),
		zap.Error(err),
	)
	if err = c.JSON(status, body); err != nil {
		log.Fatal("unrecoverable error while serializing error response", zap.Error(err))
	}
}

// ErrorHandler manages handlers erroneous returns from handlers.
// It is environment aware and will strip down the error information when used in production.
func ErrorHandler() echo.HTTPErrorHandler {
	if config.Config.IsProd() {
		return prodErrorMiddleware
	} else {
		return developmentErrorMiddleware
	}
}
