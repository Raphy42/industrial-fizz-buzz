package logger

import (
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// HttpMiddleware is a thin layer around `github.com/brpaz/echozap`, allowing us to use the application
// logger layer `zap` instead of labstack's `gommon`.
func HttpMiddleware(opts ...zap.Option) echo.MiddlewareFunc {
	return echozap.ZapLogger(New(opts...))
}
