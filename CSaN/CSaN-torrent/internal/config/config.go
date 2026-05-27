package config

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host           string `env:"TRACKER_HOST" envDefault:"0.0.0.0"`
	Port           string `env:"TRACKER_PORT" envDefault:"8080"`
	AnnounceIntSec int64  `env:"ANNOUNCE_INTERVAL_SEC" envDefault:"120"`
	PeerTimeoutSec int64  `env:"PEER_TIMEOUT_SEC" envDefault:"360"`
	DefaultNumWant int    `env:"DEFAULT_NUMWANT" envDefault:"30"`
	MaxNumWant     int    `env:"MAX_NUMWANT" envDefault:"100"`
	LogLevel       string `env:"TRACKER_LOG_LEVEL" envDefault:"info"`
	LogJSON        bool   `env:"TRACKER_LOG_JSON" envDefault:"false"`
}

func Load() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&cfg.Host, "host", cfg.Host, "server host")
	flag.StringVar(&cfg.Port, "port", cfg.Port, "server port")
	flag.Int64Var(&cfg.AnnounceIntSec, "announce-interval-sec", cfg.AnnounceIntSec, "announce interval in seconds")
	flag.Int64Var(&cfg.PeerTimeoutSec, "peer-timeout-sec", cfg.PeerTimeoutSec, "stale peer timeout in seconds")
	flag.IntVar(&cfg.DefaultNumWant, "default-numwant", cfg.DefaultNumWant, "default number of peers to return")
	flag.IntVar(&cfg.MaxNumWant, "max-numwant", cfg.MaxNumWant, "maximum number of peers to return")
	flag.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "log level: debug, info, warn, error")
	flag.BoolVar(&cfg.LogJSON, "log-json", cfg.LogJSON, "emit JSON logs to stdout")
	flag.Parse()

	if cfg.DefaultNumWant < 1 {
		cfg.DefaultNumWant = 1
	}
	if cfg.MaxNumWant < cfg.DefaultNumWant {
		cfg.MaxNumWant = cfg.DefaultNumWant
	}
	if cfg.AnnounceIntSec < 30 {
		cfg.AnnounceIntSec = 30
	}
	if cfg.PeerTimeoutSec < cfg.AnnounceIntSec {
		cfg.PeerTimeoutSec = cfg.AnnounceIntSec * 3
	}

	return cfg
}

// LogrusLevel maps config strings to logrus levels; unknown values default to info.
func (c Config) LogrusLevel() logrus.Level {
	switch strings.ToLower(strings.TrimSpace(c.LogLevel)) {
	case "debug":
		return logrus.DebugLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func (c Config) Addr() string {
	return c.Host + ":" + c.Port
}

func (c Config) AnnounceInterval() time.Duration {
	return time.Duration(c.AnnounceIntSec) * time.Second
}

func (c Config) PeerTimeout() time.Duration {
	return time.Duration(c.PeerTimeoutSec) * time.Second
}
