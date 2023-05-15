package fizzbuzz

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Raphy42/industrial-fizz-buzz/core/errors"
)

type testCase struct {
	name       string
	request    Request
	response   *Response
	statusCode int
}

var tests = []testCase{
	{
		name: "invalid int1 returns validation error",
		request: Request{
			Int1: -1,
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name: "empty str1|str2 returns validation error",
		request: Request{
			Int1:  3,
			Int2:  5,
			Limit: 10,
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name: "invalid limit returns validation error",
		request: Request{
			Int1:  69,
			Int2:  420,
			Limit: -10,
			Str1:  "Fizz",
			Str2:  "Buzz",
		},
		statusCode: http.StatusBadRequest,
	},
	{
		name: "works as expected",
		request: Request{
			Int1:  3,
			Int2:  5,
			Limit: 15,
			Str1:  "Fizz",
			Str2:  "Buzz",
		},
		response: &Response{"1", "2", "Fizz", "4", "Buzz", "Fizz", "7", "8", "Fizz", "Buzz", "11", "Fizz", "13", "14", "FizzBuzz"},
	},
	{
		name: "works as expected, even with emoji",
		request: Request{
			Int1:  3,
			Int2:  5,
			Limit: 15,
			Str1:  "😎",
			Str2:  "🤓",
		},
		response: &Response{"1", "2", "😎", "4", "🤓", "😎", "7", "8", "😎", "🤓", "11", "😎", "13", "14", "😎🤓"},
	},
}

func TestFizzBuzzHandler(t *testing.T) {
	a := assert.New(t)

	for _, test := range tests {
		if !a.True(func() bool {
			response, err := fizzBuzz(context.Background(), test.request)
			if !a.Equal(test.response, response, "responses don't match") {
				return false
			}
			if test.response != nil {
				if !a.Len(*test.response, test.request.Limit, "there are less items in the response than in the request limit") {
					return false
				}
			}
			if test.response == nil {
				if !a.Error(err, "handler returned a nil response without an error") {
					return false
				}
				httpErr, ok := err.(*errors.Error)
				if a.Truef(ok, "error is not a valid *errors.Error") {
					return a.Equal(test.statusCode, httpErr.HttpCode, "invalid http status code")
				}
				return false
			}
			return a.NoError(err, "handler returned a non nil error with a valid response")
		}(), test.name) {
			return
		}
	}
}
