package logger

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/Raphy42/industrial-fizz-buzz/core/config"
	"github.com/Raphy42/industrial-fizz-buzz/core/semconv"
)

var (
	zapOptionsCtxKey = semconv.CtxKey("zap", "options")
)

func developmentLogger(opts ...zap.Option) *zap.Logger {
	log, err := zap.NewDevelopment(opts...)
	if err != nil {
		panic(errors.Wrapf(err, "development logger initialisation failed"))
	}
	return log
}

func productionLogger(opts ...zap.Option) *zap.Logger {
	log, err := zap.NewProduction(opts...)
	if err != nil {
		panic(errors.Wrapf(err, "production logger initialisation failed"))
	}
	return log
}

// New returns a preconfigured zap.Logger using production or development config, depending on the `mode` environment
// variable.
// See also config.Config
func New(opts ...zap.Option) *zap.Logger {
	if config.Config.IsProd() {
		return productionLogger(opts...)
	}
	return developmentLogger(opts...)
}

// Inject stores zap.Option into a context.Context, to be used in conjunction with FromContext.
func Inject(ctx context.Context, opts ...zap.Option) context.Context {
	return context.WithValue(ctx, zapOptionsCtxKey, opts)
}

// FromContext fetches relevant zap.Option from within the current context.Context.
// Use Inject to store zap.Option in any given context.Context.
func FromContext(ctx context.Context) *zap.Logger {
	maybeOpts := ctx.Value(zapOptionsCtxKey)
	opts, ok := maybeOpts.([]zap.Option)
	if !ok || opts == nil {
		return New()
	}
	return New(opts...)
}
