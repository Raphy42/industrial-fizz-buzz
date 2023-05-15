package http

import (
	"context"
	"net/http"
	"path"
	"reflect"
	"runtime"

	"github.com/labstack/echo/v4"
)

type (
	// Handler is a wrapper around echo.HandlerFunc, most of the time
	Handler struct {
		operationId string
		path        string
		method      string
		impl        func(c echo.Context) error
		middlewares []echo.MiddlewareFunc
	}
	// GenericHandlerFunc is a type for generic request handlers.
	// GenericHandler expects a function of this type whenever trying to convert a generic handler func to a valid Handler.
	GenericHandlerFunc[Request any, Response any] func(ctx context.Context, request Request) (*Response, error)
	// Empty is a convenience opaque type for struct{}, use it whenever your generic handler takes no request body
	// or doesn't return a response body.
	Empty struct{}
)

// EchoHandler allows you to declare your own handlers, with full access to the echo.Context api.
func EchoHandler(route, method string, impl echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) Handler {
	fullNameOf := runtime.FuncForPC(reflect.ValueOf(impl).Pointer()).Name()
	nameOf := path.Base(fullNameOf)

	return Handler{
		operationId: nameOf,
		path:        route,
		method:      method,
		impl:        impl,
		middlewares: middlewares,
	}
}

// GenericHandler converts a generic handler into a Handler.
// This allows the user to focus on writing business code, without having to write boilerplate bind code.
// Errors are generalised and the handler will only be invoked with valid parameters.
func GenericHandler[Request any, Response any](route, method string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	fullNameOf := runtime.FuncForPC(reflect.ValueOf(impl).Pointer()).Name()
	nameOf := path.Base(fullNameOf)

	return Handler{
		operationId: nameOf,
		path:        route,
		method:      method,
		impl: func(c echo.Context) error {
			var request Request
			if err := c.Bind(&request); err != nil {
				return err
			}

			response, err := impl(c.Request().Context(), request)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, response)
		},
		middlewares: middlewares,
	}
}

func Get[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, "GET", impl, middlewares...)
}

func Post[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, "POST", impl, middlewares...)
}

func Delete[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, "DELETE", impl, middlewares...)
}

func Put[Request any, Response any](path string, impl GenericHandlerFunc[Request, Response], middlewares ...echo.MiddlewareFunc) Handler {
	return GenericHandler[Request, Response](path, "PUT", impl, middlewares...)
}
