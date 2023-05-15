package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Raphy42/industrial-fizz-buzz/core/config"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

type Server struct {
	inner *echo.Echo
}

func NewServer(handlers ...Handler) *Server {
	log := logger.New()
	e := echo.New()

	e.Debug = !config.Config.IsProd()
	e.HideBanner = true
	e.HTTPErrorHandler = ErrorHandler()
	e.Use(logger.HttpMiddleware())
	if config.Config.CorsEnabled {
		e.Use(middleware.CORS())
	}
	e.Use(middleware.RequestID())

	for _, handler := range handlers {
		log.Debug(
			"handler registered",
			zap.String("path", handler.path), zap.String("method", handler.method),
			zap.String("operationId", handler.operationId),
		)
		e.Add(handler.method, handler.path, handler.impl)
	}

	return &Server{inner: e}
}

func (s *Server) Run(ctx context.Context) error {
	log := logger.New()
	addr := config.Config.ListenAddr()

	go func() {
		log.Info("starting server", zap.String("addr", addr))
		if err := s.inner.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("ungraceful server shutdown", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	// handle parent context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-quit:
	}

	log.Info("shutting down server")

	cleanupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.inner.Shutdown(cleanupCtx); err != nil {
		return errors.Wrapf(err, "graceful server shutdown failed")
	}
	return nil
}
