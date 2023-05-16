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
	"github.com/Raphy42/industrial-fizz-buzz/core/http/metrics"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

// Server is a convenience wrapper around echo.Echo and http.Listener lifecycles
type Server struct {
	inner *echo.Echo
}

// NewServer instantiate a new Server with sane defaults, while also mounting every given Handler, and generating
// associated oas3 specs.
func NewServer(handlers ...Handler) *Server {
	log := logger.New()
	e := echo.New()

	handlers = append(handlers)

	e.Debug = !config.Config.IsProd()
	e.HideBanner = true
	e.HTTPErrorHandler = ErrorHandler()

	e.Use(middleware.RequestID())
	e.Use(logger.HttpMiddleware())
	if config.Config.CorsEnabled {
		e.Use(middleware.CORS())
	}

	oas3 := openapi()
	for _, handler := range handlers {
		if handler.reflect != nil {
			if err := handler.reflect(oas3); err != nil {
				log.Fatal("openapi3 reflection error", zap.Error(err), zap.String("path", handler.path))
			}
		}

		log.Debug(
			"handler registered",
			zap.String("path", handler.path), zap.String("method", handler.method),
			zap.String("operationId", handler.operationId),
		)
		e.Add(handler.method, handler.path, handler.impl)
	}

	schemaBytes, err := oas3.Spec.MarshalJSON()
	if err != nil {
		log.Fatal("invalid openapi3 JSON schema", zap.Error(err))
	}

	e.GET("openapi.json", func(c echo.Context) error {
		return c.Blob(http.StatusOK, "application/json", schemaBytes)
	})

	return &Server{inner: e}
}

// Run starts the server and every associated sub-systems, this call blocks and should only be called once in your application
func (s *Server) Run(ctx context.Context) error {
	log := logger.New()
	addr := config.Config.ListenAddr()

	//handlers are registered by this point
	//we can start the metrics subsystem
	metrics.StartAggregator(ctx)

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
