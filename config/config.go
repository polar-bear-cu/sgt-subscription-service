package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	GRPCPort  string
	JWTSecret string
	DB        DBConfig
}

type DBConfig struct{ Host, Port, User, Password, Name, SSLMode string }

func (db DBConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode)
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	cfg := &Config{
		Port:      env("PORT", "8080"),
		GRPCPort:  env("GRPC_PORT", "50051"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		DB: DBConfig{
			Host:     env("DB_HOST", "localhost"),
			Port:     env("DB_PORT", "5433"),
			User:     env("DB_USER", "postgres"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     env("DB_NAME", "subscriptions"),
			SSLMode:  env("DB_SSLMODE", "disable"),
		},
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET required")
	}
	return cfg, nil
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
