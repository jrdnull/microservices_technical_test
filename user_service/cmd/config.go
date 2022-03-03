package cmd

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port int `env:"PORT" envDefault:"50051"`
	DB   struct {
		User string `env:"DB_USER" envDefault:"postgres"`
		Pass string `env:"DB_PASS" envDefault:"dev"`
		Host string `env:"DB_HOST" envDefault:"db"`
		Port int    `env:"DB_PORT" envDefault:"5432"`
	}
	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:"http://jaeger:14268/api/traces"`
}

// NewConfig returns a populated Config from the environment.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	return cfg, env.Parse(cfg)
}
