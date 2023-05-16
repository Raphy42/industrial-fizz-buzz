package api

import (
	"github.com/Raphy42/industrial-fizz-buzz/api/fizzbuzz"
	"github.com/Raphy42/industrial-fizz-buzz/api/health"
	"github.com/Raphy42/industrial-fizz-buzz/api/metrics"

	"github.com/Raphy42/industrial-fizz-buzz/core/http"
)

// Handlers returns a list of http.Handler ready to be used through `server.New()`
func Handlers() []http.Handler {
	return []http.Handler{
		health.Handler,
		fizzbuzz.FizzBuzz,
		metrics.FizzBuzzMetrics,
		metrics.AllMetrics,
	}
}
