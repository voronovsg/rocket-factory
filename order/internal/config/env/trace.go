package env

import "github.com/caarlos0/env/v11"

type traceEnvConfig struct {
	CollectorEndpoint string `env:"TRACE_COLLECTOR_ENDPOINT,required"`
	ServiceName       string `env:"TRACE_SERVICE_NAME,required"`
	Environment       string `env:"TRACE_ENVIRONMENT,required"`
	ServiceVersion    string `env:"TRACE_SERVICE_VERSION,required"`
}

type traceConfig struct {
	raw traceEnvConfig
}

func NewTraceConfig() (*traceConfig, error) {
	var raw traceEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &traceConfig{raw: raw}, nil
}

func (cfg *traceConfig) CollectorEndpoint() string {
	return cfg.raw.CollectorEndpoint
}

func (cfg *traceConfig) ServiceName() string {
	return cfg.raw.ServiceName
}

func (cfg *traceConfig) Environment() string {
	return cfg.raw.Environment
}

func (cfg *traceConfig) ServiceVersion() string {
	return cfg.raw.ServiceVersion
}
