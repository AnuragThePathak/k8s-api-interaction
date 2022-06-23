package main

import (
	"log"

	"github.com/AnuragThePathak/k8s-api-interaction/server"
	"go.uber.org/zap"
	"github.com/AnuragThePathak/my-go-packages/os"
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
