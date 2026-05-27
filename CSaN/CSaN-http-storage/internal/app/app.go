package app

import (
	"context"
	"htst/internal/htst"
	"htst/pkg/config"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

func Run(config *config.Config) {
	htst := htst.New(config)
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := htst.Shutdown(ctx); err != nil {
			logrus.WithError(err).Error("Failed to shutdown HTST server")
		}
	}()

	err := htst.Start()
	if err != nil && err != http.ErrServerClosed {
		logrus.WithError(err).Error("Failed to start HTST server")
		return
	}
}
