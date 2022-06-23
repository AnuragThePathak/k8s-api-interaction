package main

import (
	"fmt"
	"log"

	"github.com/AnuragThePathak/k8s-api-interaction/endpoint"
	"github.com/AnuragThePathak/k8s-api-interaction/server"
	"go.uber.org/zap"
)

func main() {
	var logger *zap.Logger
	var err error
	{
		config := zapConfig()
		logger, err = config.Build()
		if err != nil {
			log.Fatal(fmt.Errorf("failed to initialize logger: %w", err))
		}
		defer logger.Sync()
	}
	var apiServer server.Server
	{
		config, err := serverConfig(logger)
		if err != nil {
			logger.Fatal("failed to get server config", zap.Error(err))
		}
		apiServer = server.NewServer([]server.Endpoints{
			&endpoint.PodEndpoints{},
		} ,config, logger)
	}
	apiServer.ListenAndServe()
}
