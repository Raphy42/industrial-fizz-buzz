package http

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"reflect"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/swaggest/openapi-go/openapi3"

	"github.com/Raphy42/industrial-fizz-buzz/core/http/metrics"
)

type (
	// Handler is a wrapper around echo.HandlerFunc with openapi3 and type safety in mind.
	Handler struct {
		operationId string
		path        string
		method      string
		impl        func(c echo.Context) error
		middlewares []echo.MiddlewareFunc
		reflect     func(reflector openapi3.Reflector) error
	}
	// GenericHandlerFunc is a type for generic request handlers.
	// GenericHandler expects a function of this type whenever trying to convert a generic handler func to a valid Handler.
	GenericHandlerFunc[Request any, Response any] func(ctx context.Context, request Request) (*Response, error)
	// Empty is a convenience opaque type for struct{}, use it whenever your generic handler takes no request body
	// or doesn't return a response body.
	Empty struct{}
)

// EchoHandler allows you to declare your own handlers, with full access to the echo.Context api.
// This is currently unused and not connected with openapi3 systems.
func EchoHandler(route, method string, impl echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) Handler {
	fullNameOf := runtime.FuncForPC(reflect.ValueOf(impl).Pointer()).Name()
	nameOf := path.Base(fullNameOf)

	// todo implement when needed
	// requestChan := metrics.NewRequestCounter(route)

	h := Handler{
		operationId: nameOf,
		path:        route,
		method:      method,
		impl:        impl,
		middlewares: middlewares,
	}
	return h
}

// GenericHandler converts a generic handler into a Handler.
// This allows the user to focus on writing business code, without having to write boilerplate bind code.
// Errors are generalised and the handler will only be invoked with valid parameters.
func GenericHandler[Request any, Response any](route, method string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	fullNameOf := runtime.FuncForPC(reflect.ValueOf(impl).Pointer()).Name()
	nameOf := path.Base(fullNameOf)
	requestChan := metrics.NewRequestCounter(route)

	h := Handler{
		operationId: nameOf,
		path:        route,
		method:      method,
		impl: func(c echo.Context) error {
			var request Request
			if err := c.Bind(&request); err != nil {
				return err
			}

			// we serialize the complete request object to JSON
			// this is different from a JSON request body, as echo allows query, path, and json parameters
			// through tag reflection
			ctx := c.Request().Context()
			go func() {
				buf, err := json.Marshal(request)
				if err != nil {
					panic(err)
				}
				select {
				case <-ctx.Done():
					return
				case requestChan <- buf:
					return
				}
			}()

			response, err := impl(ctx, request)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, response)
		},
		middlewares: middlewares,
	}
	fn := func(reflector openapi3.Reflector) error {
		return registerOperation[Request, Response](reflector, h)
	}
	h.reflect = fn
	return h
}

// Get is a convenience wrapper around GenericHandler
func Get[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, http.MethodGet, impl, middlewares...)
}

// Post is a convenience wrapper around GenericHandler
func Post[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, http.MethodPost, impl, middlewares...)
}

// Delete is a convenience wrapper around GenericHandler
func Delete[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, http.MethodDelete, impl, middlewares...)
}

// Put is a convenience wrapper around GenericHandler
func Put[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, http.MethodPut, impl, middlewares...)
}
