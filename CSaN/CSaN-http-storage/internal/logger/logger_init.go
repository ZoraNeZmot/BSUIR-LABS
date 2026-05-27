package logger

import (
	"htst/pkg/config"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var logFile *os.File

func Init() func() {
	config := config.GetConfig()
	if config.App.Logger.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{})
	}
	level, err := logrus.ParseLevel(config.App.Logger.Level)
	if err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}
	logFile, err = os.OpenFile(config.App.Logger.File, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err == nil {
		logrus.SetOutput(io.MultiWriter(logFile, os.Stdout))
	}
	logrus.SetReportCaller(config.IS_DEBUG)
	return func() {
		if logFile != nil {
			err := logFile.Close()
			if err != nil {
				logrus.WithError(err).Error("Failed to close log file")
			}
		}
	}
}
