package api

import (
	"github.com/Raphy42/industrial-fizz-buzz/api/fizzbuzz"
	"github.com/Raphy42/industrial-fizz-buzz/api/health"

	"github.com/Raphy42/industrial-fizz-buzz/core/http"
)

func Handlers() []http.Handler {
	return []http.Handler{
		health.Handler,
		fizzbuzz.FizzBuzz,
	}
}
