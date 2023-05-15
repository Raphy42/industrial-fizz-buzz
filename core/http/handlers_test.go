package http

import (
	"context"

	"go.uber.org/zap"

	"github.com/Raphy42/industrial-fizz-buzz/core/errors"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

func ExampleGenericHandler() {
	// package users

	// request type, see `https://echo.labstack.com/guide/binding/` for more information
	type Request struct {
		UserId string `param:"id"`
	}
	// response type, aka what is returned in case the handler is successful
	type Response struct {
		Username string `json:"username"`
	}

	var (
		// non thread-safe in memory database stub
		users = map[string]string{
			"914358c7-bb4b-49eb-9d8a-7044974d5437": "alice",
			"43523879-f9e8-42b1-be8e-27ace4c2fc3e": "bob",
		}
		// the handler implementation, will commonly be its own function instead of an anonymous closure.
		// use http.Empty if your handler takes no request body, and *http.Empty if it returns an empty body.
		userHandler = func(ctx context.Context, request Request) (*Response, error) {
			log := logger.FromContext(ctx)
			log.Debug("received request", zap.String("user.id", request.UserId))

			username, ok := users[request.UserId]
			if ok {
				return nil, errors.NotFound()
			}

			return &Response{Username: username}, nil
		}
		// the final handler, this should be exported, and used by server.New whenever needed.
		// you can also add optional middlewares as the end of the call, and they will be only mounted for this
		// particular route.
		UserHandler = GenericHandler("/api/v1/users/:id", "GET", userHandler)
	)

	// final usage
	NewServer(UserHandler)
}
