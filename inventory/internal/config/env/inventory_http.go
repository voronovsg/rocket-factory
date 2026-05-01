package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvHTTPConfig struct {
	Host string `env:"HTTP_HOST,required"`
	Port string `env:"HTTP_PORT,required"`
}

type inventoryHTTPConfig struct {
	raw inventoryEnvHTTPConfig
}

func NewInventoryHTTPConfig() (*inventoryHTTPConfig, error) {
	var raw inventoryEnvHTTPConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryHTTPConfig{raw: raw}, nil
}

func (cfg *inventoryHTTPConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
