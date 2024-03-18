package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Log        `yaml:"logger"`
		HTTPServer `yaml:"http_server"`
		StorageConfig
	}

	Log struct {
		Env string `env-required:"true" yaml:"env" env:"LOG_LEVEL"`
	}

	HTTPServer struct {
		Address     string        `env-required:"true" yaml:"address" env:"HTTP_ADDRESS"`
		Timeout     time.Duration `env-required:"true" yaml:"timeout" env:"HTTP_TIMEOUT"`
		IdleTimeout time.Duration `env-required:"true" yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT"`
		User string `env-required:"true" yaml:"user" env:"HTTP_USER"`
		Pass string `env-required:"true" yaml:"pass" env:"HTTP_PASS"`
	}

	StorageConfig struct {
		URL string `env-required:"true" env:"PG_URL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
