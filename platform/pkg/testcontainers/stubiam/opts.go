package stubiam

import (
	"io"

	"github.com/testcontainers/testcontainers-go/wait"
)

type Option func(*Config)

func WithNetworkName(network string) Option {
	return func(c *Config) {
		c.NetworkName = network
	}
}

func WithContainerName(name string) Option {
	return func(c *Config) {
		c.ContainerName = name
	}
}

func WithPort(port string) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithLogOutput(out io.Writer) Option {
	return func(c *Config) {
		c.LogOutput = out
	}
}

func WithStartupWait(strategy wait.Strategy) Option {
	return func(c *Config) {
		c.StartupWait = strategy
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}
