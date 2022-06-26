package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/AnuragThePathak/k8s-api-interaction/server"
	"github.com/AnuragThePathak/my-go-packages/os"
	"go.uber.org/zap"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func serverConfig(logger *zap.Logger) (server.ServerConfig, error) {
	config := server.ServerConfig{}
	var err error
	config.Port, err = os.GetEnvAsInt("PORT", 8080)
	if err != nil {
		return config, err
	}
	logger.Debug("sever config:", zap.Int("PORT", config.Port))
	return config, nil
}

func zapConfig() zap.Config {
	env, _ := os.GetEnv("ENV")
	log.Println("ENV:", env)
	if env == "production" {
		return zap.NewProductionConfig()
	}
	return zap.NewDevelopmentConfig()
}

func kubeConfig() (*rest.Config, error) {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		log.Println(home)
		kubeconfig = flag.String("kube-config", filepath.Join(home, ".kube",
		"config"), "(optional) absolute path to the kubeconfig file")
		} else {
		kubeconfig = flag.String("kube-config", "",
			"absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		return config, fmt.Errorf("failed to get kubeconfig: %w", err)
	}
	return config, nil
}
