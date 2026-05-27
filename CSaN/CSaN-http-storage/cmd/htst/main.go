package main

import (
	"htst/internal/app"
	"htst/internal/logger"
	"htst/pkg/config"

	"github.com/sirupsen/logrus"
)

func main() {
	config := config.GetConfig()
	cleanupLogger := logger.Init()
	if config.IS_DEBUG {
		logrus.Info("Config: ", *config)
	}
	defer cleanupLogger()
	app.Run(config)
}
