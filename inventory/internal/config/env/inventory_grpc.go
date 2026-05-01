package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type inventoryEnvGRPCConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw inventoryEnvGRPCConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryEnvGRPCConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
