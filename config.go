package main

import (
	"os"
	"strconv"

	"go.uber.org/zap"
)

func serverConfig() (ServerConfig, error) {
	config := ServerConfig{}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}
	var err error
	config.Port, err = strconv.Atoi(port)
	if err != nil {
		return config, err
	}
	return config, nil
}

func zapConfig() zap.Config {
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "dev"
	}
	if env == "production" {
		return zap.NewProductionConfig()
	}
	return zap.NewDevelopmentConfig()
}
