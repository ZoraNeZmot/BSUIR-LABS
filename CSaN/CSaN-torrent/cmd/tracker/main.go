package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"torrent/internal/config"
	trackerhttp "torrent/internal/http"
	"torrent/internal/tracker"

	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.Load()

	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(cfg.LogrusLevel())
	if cfg.LogJSON {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{})
	}
	log.Info(cfg)
	store := tracker.NewStore(time.Now().UnixNano())
	service := tracker.NewService(store, cfg.AnnounceIntSec, cfg.DefaultNumWant, cfg.MaxNumWant)

	mux := http.NewServeMux()
	handlers := trackerhttp.NewHandlers(log, service, store)
	handlers.Register(mux)

	server := &http.Server{
		Addr:              cfg.Addr(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go tracker.StartCleanup(ctx, log, store, cfg.PeerTimeout())

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.WithError(err).Error("server shutdown error")
		}
	}()

	log.WithFields(logrus.Fields{
		"addr":                  server.Addr,
		"announce_interval_sec": cfg.AnnounceIntSec,
		"peer_timeout_sec":      cfg.PeerTimeoutSec,
		"log_level":             cfg.LogLevel,
		"log_json":              cfg.LogJSON,
	}).Info("tracker listening")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("server error")
		os.Exit(1)
	}
	log.Info("tracker stopped")
}
