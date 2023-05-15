package errors

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/pkg/errors"
)

type (
	Trace struct {
		File   string `json:"file"`
		Line   int    `json:"line"`
		Caller string `json:"caller"`
	}
	// Error is a specialised error type, which contains important metadata.
	// This type is used by the http.ErrorHandler, so it should be used in handlers whenever possible.
	Error struct {
		error    `json:"error"`
		HttpCode int    `json:"status"`
		Message  string `json:"message"`
		Trace    *Trace `json:"trace,omitempty"`
	}
)

func newError(skipCallers int, parent error, statusCode int, format string, args ...any) *Error {
	var trace *Trace
	var err error

	caller, file, line, ok := runtime.Caller(skipCallers + 1)
	msg := fmt.Sprintf(format, args...)
	if ok {
		desc := runtime.FuncForPC(caller)
		if desc != nil {
			trace = &Trace{
				File:   file,
				Line:   line,
				Caller: desc.Name(),
			}
		}
	}
	if parent != nil {
		err = errors.Wrapf(parent, msg)
	} else {
		err = errors.Errorf(msg)
	}

	return &Error{
		error:    err,
		Message:  msg,
		HttpCode: statusCode,
		Trace:    trace,
	}
}

// NotImplemented is a placeholder error which can be used whenever needed to mark a work in progress.
// StatusCode: 501
func NotImplemented() error {
	return newError(1, nil, http.StatusNotImplemented, "this part has not been yet implemented")
}

// BadRequest wraps an optional error and a message with args, into an error containing important metadata for
// downstream processing.
// StatusCode: 400
func BadRequest(err error, format string, args ...any) error {
	return newError(1, err, http.StatusBadRequest, format, args...)
}

// NotFound is a convenience wrapper for not found resource error.
// StatusCode: 404
func NotFound() error {
	return newError(1, nil, http.StatusNotFound, "resource not found")
}
