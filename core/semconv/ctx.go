package semconv

import "strings"

const namespace = "fizzbuzz"

// CtxKey returns an identifier which can be used to inject value in context.Context
func CtxKey(fields ...string) string {
	fields = append([]string{namespace}, fields...)
	return strings.Join(fields, ".")
}
