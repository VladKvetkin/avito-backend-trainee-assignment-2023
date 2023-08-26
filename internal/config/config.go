package config

import (
	"net/url"

	"github.com/caarlos0/env"
)

type Config struct {
	Address     string `env:"ADDRESS"`
	DatabaseDSN string `env:"DATABASE_DSN"`
}

func NewConfig() (Config, error) {
	config := Config{}

	if err := env.Parse(&config); err != nil {
		return Config{}, err
	}

	if err := config.validateConfig(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c *Config) validateConfig() error {
	_, err := url.ParseRequestURI(c.Address)
	if err != nil {
		return err
	}

	return nil
}
