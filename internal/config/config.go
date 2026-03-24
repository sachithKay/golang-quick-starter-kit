package config

import (
	"log/slog"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string `env:"PORT" envDefault:"8080"`
	DB   struct {
		URL             string        `env:"DB_URL,required"` // Fails if missing
		MaxConnections  int           `env:"DB_MAX_CONNS" envDefault:"20"`
		MinConnections  int           `env:"DB_MIN_CONNS" envDefault:"5"`
		ConnMaxIdleTime time.Duration `env:"DB_CONN_IDLE_TIME" envDefault:"30m"`
	}
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		slog.Error("Failed to parse config", "error", err)
		return nil, err
	}

	slog.Info("config loaded successfully", "config", cfg)
	return cfg, nil
}
