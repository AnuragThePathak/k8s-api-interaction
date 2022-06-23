package main

import (
	"fmt"
	"log"

	"github.com/AnuragThePathak/k8s-api-interaction/endpoint"
	"github.com/AnuragThePathak/k8s-api-interaction/server"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
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
	var clientSet *kubernetes.Clientset
	{
		config, err := kubeConfig()
		if err != nil {
			logger.Panic("failed to get kube config", zap.Error(err))
		}
		clientSet = kubernetes.NewForConfigOrDie(config)
	}
	var apiServer server.Server
	{
		config, err := serverConfig(logger)
		if err != nil {
			logger.Fatal("failed to get server config", zap.Error(err))
		}
		apiServer = server.NewServer([]server.Endpoints{
			&endpoint.PodEndpoints{
				ClientSet: clientSet,
			},
			&endpoint.ServiceEndpoints{
				ClientSet: clientSet,
			},
		} ,config, logger)
	}
	apiServer.ListenAndServe()
}
