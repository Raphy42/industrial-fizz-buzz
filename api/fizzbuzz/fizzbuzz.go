package fizzbuzz

import (
	"context"
	"strconv"

	"go.uber.org/zap"

	"github.com/Raphy42/industrial-fizz-buzz/core/config"
	"github.com/Raphy42/industrial-fizz-buzz/core/errors"
	"github.com/Raphy42/industrial-fizz-buzz/core/http"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

type (
	Request struct {
		Int1  int    `query:"int1"`
		Int2  int    `query:"int2"`
		Limit int    `query:"limit"`
		Str1  string `query:"str1"`
		Str2  string `query:"str2"`
	}
	Response []string
)

var FizzBuzz = http.GenericHandler("/api/v1/fizzbuzz", "GET", fizzBuzz)

func fizzBuzzImpl(number int, request Request) string {
	switch 0 {
	case number % (request.Int1 * request.Int2):
		return request.Str1 + request.Str2
	case number % request.Int1:
		return request.Str1
	case number % request.Int2:
		return request.Str2
	default:
		return strconv.FormatInt(int64(number), 10)
	}
}

func fizzBuzz(ctx context.Context, request Request) (*Response, error) {
	log := logger.FromContext(ctx)

	log.Debug("new fizzbuzz request",
		zap.Strings("words", []string{request.Str1, request.Str2}),
		zap.Ints("ints", []int{request.Int1, request.Int2}),
		zap.Int("limit", request.Limit),
	)

	if request.Limit < 0 {
		return nil, errors.BadRequest(nil, "`limit` query parameter cannot be negative")
	}
	if request.Int1 <= 0 || request.Int2 <= 0 {
		return nil, errors.BadRequest(nil, "both `int1` and `int2` query parameters must be valid positive non-zero integer")
	}

	if !config.Config.FizzBuzzAllowEmptyStr && (request.Str1 == "" || request.Str2 == "") {
		msg := "both `str1` and `str2` query parameters must be set, empty words have been disallowed through configuration"
		return nil, errors.BadRequest(nil, msg)
	}

	results := make(Response, request.Limit)
	for i := 1; i < request.Limit+1; i++ {
		results[i-1] = fizzBuzzImpl(i, request)
	}

	return &results, nil
}