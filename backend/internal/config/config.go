package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL       string `env:"DATABASE_URL,required"`
	SupabaseURL       string `env:"SUPABASE_URL,required"`
	SupabaseJWTSecret string `env:"SUPABASE_JWT_SECRET,required"`
	Port              string `env:"PORT" envDefault:"8080"`
	Environment       string `env:"ENVIRONMENT" envDefault:"development"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
