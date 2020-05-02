package main

import (
	"os"
	"strings"

	"go.uber.org/zap"
)

func getLogger(module string) *zap.SugaredLogger {
	logLevel := os.Getenv("LOG_LEVEL")
	upperModule := strings.ToUpper(module)
	if os.Getenv("LOG_LEVEL_"+upperModule) != "" {
		logLevel = os.Getenv("LOG_LEVEL_" + upperModule)
	}

	runEnv := os.Getenv("RUN_ENV")
	var config zap.Config
	if strings.ToUpper(runEnv) == "PROD" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.Level.UnmarshalText([]byte(logLevel))
	log, _ := config.Build()

	return log.Named(module).Sugar()
}
