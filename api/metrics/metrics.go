package metrics

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/Raphy42/industrial-fizz-buzz/core/http"
	"github.com/Raphy42/industrial-fizz-buzz/core/http/metrics"
)

type (
	// Response used by both the AllMetrics and FizzBuzzMetrics endpoints
	Response struct {
		Request map[string]any `json:"request"`
		Hits    uint           `json:"hits"`
	}
)

var (
	// FizzBuzzMetrics handles GET /api/v1/metrics/request/fizzbuzz
	FizzBuzzMetrics = http.Get("/api/v1/metrics/request/fizzbuzz", fizzbuzzMetrics)
	// AllMetrics handles GET /api/v1/metrics/request
	AllMetrics = http.Get("/api/v1/metrics/request", listMetrics)
)

func listMetrics(_ context.Context, _ http.Empty) (*map[string]Response, error) {
	responses := make(map[string]Response)
	result, err := metrics.Top()
	for route, top := range result {
		var req map[string]any
		if top.Bytes != nil {
			if err = json.Unmarshal(top.Bytes, &req); err != nil {
				return nil, errors.Wrapf(err, "json unmarshaling of top request body for '%s' failed", top.Route)
			}
		}
		responses[route] = Response{
			Request: req,
			Hits:    top.Hits,
		}
	}
	return &responses, nil
}

func fizzbuzzMetrics(_ context.Context, _ http.Empty) (*Response, error) {
	const endpoint = "/api/v1/fizzbuzz"

	result, err := metrics.Top(endpoint)
	if err != nil {
		return nil, err
	}
	top := result[endpoint]
	var req map[string]any
	if top.Bytes != nil {
		if err = json.Unmarshal(top.Bytes, &req); err != nil {
			return nil, errors.Wrapf(err, "json unmarshaling of top request body for '%s' failed", top.Route)
		}
	}
	return &Response{
		Request: req,
		Hits:    top.Hits,
	}, nil
}
