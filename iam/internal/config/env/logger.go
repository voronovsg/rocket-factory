package env

import (
	"github.com/caarlos0/env/v11"
)

type loggerConfig struct {
	LoggerLevel  string `env:"LOGGER_LEVEL,required"`
	LoggerAsJson bool   `env:"LOGGER_AS_JSON,required"`
}

func NewLoggerConfig() (*loggerConfig, error) {
	var cfg loggerConfig
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.LoggerLevel
}

func (cfg *loggerConfig) AsJson() bool {
	return cfg.LoggerAsJson
}
