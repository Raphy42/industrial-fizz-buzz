package main

import (
	"context"

	"go.uber.org/zap"

	"github.com/Raphy42/industrial-fizz-buzz/api"
	"github.com/Raphy42/industrial-fizz-buzz/core/http"
	"github.com/Raphy42/industrial-fizz-buzz/core/logger"
)

func main() {
	log := logger.New()
	ctx := context.Background()

	server := http.NewServer(api.Handlers()...)

	if err := server.Run(ctx); err != nil {
		log.Fatal("server crashed", zap.Error(err))
	}
}
