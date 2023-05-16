package metrics

import "context"

// StartAggregator runs the internal channel multiplexing and thread-safe state management of the globalRegistry
func StartAggregator(ctx context.Context) {
	globalRegistry.start(ctx)
}

// NewRequestCounter allocates a new request counter, this is used internally by the handler wrapper to extract metrics.
func NewRequestCounter(path string) chan<- []byte {
	return globalRegistry.newRequestCounter(path)
}

// Top returns a map containing all top requests by routes, if the input is empty, or each route and its associated
// top request count.
func Top(routes ...string) (map[string]TopRequest, error) {
	return globalRegistry.top(routes...)
}

// TopRequest represents the current top request for a given route
type TopRequest struct {
	Route string
	Bytes []byte
	Hits  uint
}
