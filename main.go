package main

import (
	"log"

	"go.uber.org/zap"
)

func main() {
	var logger *zap.Logger
	var err error
	{
		config := zapConfig()
		logger, err = config.Build()
		if err != nil {
			log.Fatal(err)
		}
		defer logger.Sync()
	}
}
