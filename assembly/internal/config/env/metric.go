package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type metricServerEnvConfig struct {
	CollectorEndpoint string        `env:"METRIC_COLLECTOR_ENDPOINT,required"`
	CollectorInterval time.Duration `env:"METRIC_COLLECTOR_INTERVAL,required"`
}

type metricServerConfig struct {
	raw metricServerEnvConfig
}

func NewMetricServerConfig() (*metricServerConfig, error) {
	var raw metricServerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &metricServerConfig{raw: raw}, nil
}

func (cfg *metricServerConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *metricServerConfig) CollectorInterval() time.Duration {
	return cfg.raw.CollectorInterval
}
