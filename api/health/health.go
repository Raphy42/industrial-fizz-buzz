package health

import (
	"context"

	"github.com/Raphy42/industrial-fizz-buzz/core/http"
)

var (
	// Handler handles GET /health
	Handler = http.Get("/health", healthCheck)
)

func healthCheck(_ context.Context, _ http.Empty) (*http.Empty, error) {
	return &http.Empty{}, nil
}
