package config

import (
	"sync"

	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	IS_DEBUG bool `env:"IS_DEBUG" envDefault:"false"`
	App      AppConfig
}
type AppConfig struct {
	Listen ListenConfig `envPrefix:"LISTEN_"`
	Logger LoggerConfig `envPrefix:"LOGGER_"`
	HTST   HTSTConfig   `envPrefix:"HTST_"`
}
type HTSTConfig struct {
	Root string `env:"ROOT" envDefault:"./storage"`
	// TODO: Index file handling
}
type ListenConfig struct {
	Port string `env:"PORT" envDefault:"8080"`
	Host string `env:"HOST" envDefault:"0.0.0.0"`
}
type LoggerConfig struct {
	File   string `env:"FILE" envDefault:"htst.log"`
	Level  string `env:"LEVEL" envDefault:"info"`  // trace, debug, info, warn, error, fatal, panic
	Format string `env:"FORMAT" envDefault:"json"` // text, json
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{}
		env.Parse(instance)

	})
	return instance
}
