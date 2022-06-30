package main

import (
	"fmt"
	"log"

	"github.com/AnuragThePathak/k8s-api-interaction/endpoint"
	"github.com/AnuragThePathak/k8s-api-interaction/server"
	openelb "github.com/openelb/openelb/api/v1alpha2"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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
		kubeConfig, err := kubeConfig()
		if err != nil {
			logger.Panic("failed to get kube config", zap.Error(err))
		}
		if clientSet, err = kubernetes.NewForConfig(kubeConfig); err != nil {
			logger.Panic("failed to get kube client", zap.Error(err))
		}
	}
	var runtimeclientSet client.Client
	{
		crScheme := runtime.NewScheme()
		if err = openelb.AddToScheme(crScheme); err != nil {
			logger.Panic("failed to add openelb to scheme", zap.Error(err))
		}
		if runtimeclientSet, err = client.New(config.GetConfigOrDie(),
			client.Options{
				Scheme: crScheme,
			}); err != nil {
			logger.Panic("failed to create runtime client", zap.Error(err))
		}
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
			&endpoint.OpenelbEndpoints{
				Client: runtimeclientSet,
			},
		}, config, logger)
	}
	apiServer.ListenAndServe()
}
