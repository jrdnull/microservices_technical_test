package cmd

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
	DB   struct {
		User string `env:"DB_USER" envDefault:"postgres"`
		Pass string `env:"DB_PASS" envDefault:"dev"`
		Host string `env:"DB_HOST" envDefault:"db"`
		Port int    `env:"DB_PORT" envDefault:"5432"`
	}
	UserService struct {
		Address string `env:"USER_SERVICE_ADDRESS" envDefault:"user_service:50051"`
	}
	JaegerEndpoint string `env:"JAEGER_ENDPOINT" envDefault:"http://jaeger:14268/api/traces"`
	// Specifies the required amount of tag ids to provide to the
	// /articles tag filter.
	ArticleTagFilterInputs int `env:"ARTICLE_TAG_FILTER_INPUTS" envDefault:"2"`
}

// NewConfig returns a populated Config from the environment.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	return cfg, env.Parse(cfg)
}
