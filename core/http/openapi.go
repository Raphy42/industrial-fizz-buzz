package http

import (
	"net/http"

	"github.com/swaggest/openapi-go/openapi3"
)

func openapi() openapi3.Reflector {
	reflector := openapi3.Reflector{}
	reflector.Spec = &openapi3.Spec{
		Openapi: "3.0.3",
		Servers: []openapi3.Server{
			{
				URL: "http://localhost:8080",
			},
		},
		Info: openapi3.Info{
			Title:   "FizzBuzz api",
			Version: "v0.1.0",
		},
	}
	return reflector
}

func registerOperation[Request any, Response any](reflector openapi3.Reflector, h Handler) error {
	op := openapi3.Operation{}
	if err := reflector.SetRequest(&op, new(Request), h.method); err != nil {
		return err
	}
	if err := reflector.SetJSONResponse(&op, new(Response), http.StatusOK); err != nil {
		return err
	}
	return reflector.Spec.AddOperation(h.method, h.path, op)
}
